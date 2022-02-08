FROM golang:1.17-alpine AS builder

WORKDIR /build
COPY . /build/

RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:3.15

COPY --from=builder /build/mailbowl /usr/bin/
COPY ./config.example.yml /tmp/

RUN addgroup -S mailbowl \
 && adduser -S -D -H -G mailbowl mailbowl \
 && mkdir -p /etc/mailbowl \
 && mv /tmp/config.example.yml /etc/mailbowl/mailbowl.yaml \
 && chown -R mailbowl:mailbowl /etc/mailbowl

USER mailbowl:mailbowl

EXPOSE 10465 10587

CMD ["/usr/bin/mailbowl", "--config", "/etc/mailbowl/mailbowl.yaml"]
