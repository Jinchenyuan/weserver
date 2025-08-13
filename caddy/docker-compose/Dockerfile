FROM caddy:builder AS builder

ENV GOPROXY=https://goproxy.cn,direct
RUN xcaddy build \
    --with github.com/caddy-dns/cloudflare

FROM caddy:latest

COPY --from=builder /usr/bin/caddy /usr/bin/caddy