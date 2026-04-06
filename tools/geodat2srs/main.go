package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sagernet/sing-box/common/geosite"
	"github.com/sagernet/sing-box/common/srs"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

type source struct {
	Name string
	Kind string
	URLs []string
}

type sourceFlags []source

func (s *sourceFlags) String() string {
	return fmt.Sprint([]source(*s))
}

func (s *sourceFlags) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("source must be in format kind:name=url1,url2")
	}

	left := strings.SplitN(parts[0], ":", 2)
	if len(left) != 2 {
		return fmt.Errorf("source must be in format kind:name=url1,url2")
	}

	urls := make([]string, 0)
	for _, rawURL := range strings.Split(parts[1], ",") {
		url := strings.TrimSpace(rawURL)
		if url != "" {
			urls = append(urls, url)
		}
	}

	if len(urls) == 0 {
		return fmt.Errorf("source must contain at least one url")
	}

	*s = append(*s, source{
		Kind: strings.TrimSpace(left[0]),
		Name: strings.TrimSpace(left[1]),
		URLs: urls,
	})

	return nil
}

func main() {
	outputDir := flag.String("output-dir", "rules", "Output directory")
	var sources sourceFlags
	flag.Var(&sources, "source", "Source in format kind:name=url1,url2")
	flag.Parse()

	if len(sources) == 0 {
		fail(fmt.Errorf("at least one -source argument is required"))
	}

	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		fail(err)
	}

	for _, source := range sources {
		input, usedURL, err := readFirstAvailable(source.URLs)
		if err != nil {
			fail(fmt.Errorf("download %s: %w", source.Name, err))
		}

		fmt.Fprintf(os.Stdout, "using %s from %s\n", source.Name, usedURL)

		switch source.Kind {
		case "geoip":
			err = writeGeoIP(filepath.Join(*outputDir, source.Name+".srs"), input)
		case "geosite":
			err = writeGeoSiteCategories(*outputDir, source.Name, input)
		default:
			err = fmt.Errorf("unsupported source kind: %s", source.Kind)
		}

		if err != nil {
			fail(fmt.Errorf("write %s: %w", source.Name, err))
		}
	}
}

func writeGeoIP(outputPath string, input []byte) error {
	var geoIPList routercommon.GeoIPList
	if err := proto.Unmarshal(input, &geoIPList); err != nil {
		return fmt.Errorf("parse geoip.dat: %w", err)
	}

	cidrs := make([]string, 0)
	seen := make(map[string]struct{})
	for _, geoIP := range geoIPList.Entry {
		for _, cidr := range geoIP.Cidr {
			ip := net.IP(cidr.GetIp())
			if len(ip) == 0 {
				continue
			}
			value := ip.String() + "/" + fmt.Sprint(cidr.GetPrefix())
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}
			cidrs = append(cidrs, value)
		}
	}

	return writeRuleSet(outputPath, option.PlainRuleSet{
		Rules: []option.HeadlessRule{
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultHeadlessRule{
					IPCIDR: cidrs,
				},
			},
		},
	})
}

func writeGeoSiteCategories(outputDir string, prefix string, input []byte) error {
	var geoSiteList routercommon.GeoSiteList
	if err := proto.Unmarshal(input, &geoSiteList); err != nil {
		return fmt.Errorf("parse geosite.dat: %w", err)
	}

	for _, entry := range geoSiteList.Entry {
		code := normalizeName(entry.GetCode())
		if code == "" {
			code = normalizeName(entry.GetCountryCode())
		}
		if code == "" {
			continue
		}

		items := make([]geosite.Item, 0)
		seen := make(map[string]struct{})
		for _, domain := range entry.Domain {
			switch domain.Type {
			case routercommon.Domain_Plain:
				appendItem(&items, seen, geosite.RuleTypeDomainKeyword, domain.Value)
			case routercommon.Domain_Regex:
				appendItem(&items, seen, geosite.RuleTypeDomainRegex, domain.Value)
			case routercommon.Domain_RootDomain:
				if strings.Contains(domain.Value, ".") {
					appendItem(&items, seen, geosite.RuleTypeDomain, domain.Value)
				}
				appendItem(&items, seen, geosite.RuleTypeDomainSuffix, "."+domain.Value)
			case routercommon.Domain_Full:
				appendItem(&items, seen, geosite.RuleTypeDomain, domain.Value)
			}
		}

		if len(items) == 0 {
			continue
		}

		compiled := geosite.Compile(items)
		outputPath := filepath.Join(outputDir, prefix+"-"+code+".srs")
		err := writeRuleSet(outputPath, option.PlainRuleSet{
			Rules: []option.HeadlessRule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultHeadlessRule{
						Domain:        compiled.Domain,
						DomainSuffix:  compiled.DomainSuffix,
						DomainKeyword: compiled.DomainKeyword,
						DomainRegex:   compiled.DomainRegex,
					},
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func appendItem(items *[]geosite.Item, seen map[string]struct{}, ruleType geosite.ItemType, value string) {
	key := fmt.Sprintf("%d:%s", ruleType, value)
	if _, ok := seen[key]; ok {
		return
	}

	seen[key] = struct{}{}
	*items = append(*items, geosite.Item{
		Type:  ruleType,
		Value: value,
	})
}

func writeRuleSet(outputPath string, ruleSet option.PlainRuleSet) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return srs.Write(outputFile, ruleSet)
}

func normalizeName(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, ":", "-")
	value = strings.ReplaceAll(value, "/", "-")
	value = strings.ReplaceAll(value, "\\", "-")
	value = strings.ReplaceAll(value, " ", "-")
	return value
}

func readFirstAvailable(paths []string) ([]byte, string, error) {
	var errors []string
	for _, path := range paths {
		data, err := readFile(path)
		if err == nil {
			return data, path, nil
		}
		errors = append(errors, err.Error())
	}

	return nil, "", fmt.Errorf(strings.Join(errors, "; "))
}

func readFile(path string) ([]byte, error) {
	switch {
	case strings.HasPrefix(strings.ToLower(path), "http://"), strings.HasPrefix(strings.ToLower(path), "https://"):
		resp, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: http status code %d", path, resp.StatusCode)
		}

		return io.ReadAll(resp.Body)
	default:
		return os.ReadFile(path)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
