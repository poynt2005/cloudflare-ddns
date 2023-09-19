FROM golang:alpine as go-builder
WORKDIR /progress
WORKDIR /usr/local/go/src/
COPY ddns ./ddns
RUN cd ddns \
    && go build -o cloudflare_ddns main.go \
    && chmod +x cloudflare_ddns \
    && mv cloudflare_ddns /progress/cloudflare_ddns
RUN apk --update add --no-cache ca-certificates openssl \
    && update-ca-certificates

FROM scratch
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /progress/cloudflare_ddns /cloudflare_ddns
CMD ["/cloudflare_ddns"]
