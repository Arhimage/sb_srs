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

func main() {
	geoIPPath := flag.String("geoip", "geoip.dat", "Path or URL to geoip.dat")
	geoSitePath := flag.String("geosite", "geosite.dat", "Path or URL to geosite.dat")
	outputDir := flag.String("output-dir", "rules", "Output directory")
	flag.Parse()

	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		fail(err)
	}

	geoIPBytes, err := readFile(*geoIPPath)
	if err != nil {
		fail(err)
	}

	geoSiteBytes, err := readFile(*geoSitePath)
	if err != nil {
		fail(err)
	}

	if err := writeGeoIP(filepath.Join(*outputDir, "geoip.srs"), geoIPBytes); err != nil {
		fail(err)
	}

	if err := writeGeoSite(filepath.Join(*outputDir, "geosite.srs"), geoSiteBytes); err != nil {
		fail(err)
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

func writeGeoSite(outputPath string, input []byte) error {
	var geoSiteList routercommon.GeoSiteList
	if err := proto.Unmarshal(input, &geoSiteList); err != nil {
		return fmt.Errorf("parse geosite.dat: %w", err)
	}

	items := make([]geosite.Item, 0)
	seen := make(map[string]struct{})
	for _, entry := range geoSiteList.Entry {
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
	}

	compiled := geosite.Compile(items)
	return writeRuleSet(outputPath, option.PlainRuleSet{
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

func readFile(path string) ([]byte, error) {
	switch {
	case strings.HasPrefix(strings.ToLower(path), "http://"), strings.HasPrefix(strings.ToLower(path), "https://"):
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get remote file %s, http status code %d", path, resp.StatusCode)
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
