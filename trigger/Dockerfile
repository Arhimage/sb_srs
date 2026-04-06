FROM alpine:3.19

RUN apk add --no-cache curl

COPY trigger.sh /trigger.sh
RUN chmod +x /trigger.sh

ENTRYPOINT ["/trigger.sh"]
