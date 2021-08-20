# From https://tailscale.com/kb/1108/cloudrun/
FROM alpine:latest as tailscale
WORKDIR /app
COPY . ./
ENV TSFILE=tailscale_1.12.3_amd64.tgz
RUN wget https://pkgs.tailscale.com/stable/${TSFILE} && \
  tar xzf ${TSFILE} --strip-components=1
COPY . ./

FROM caddy:2.4.3-alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Copy binary to production image
COPY package/caddy-tailscale/start.sh /app/start.sh
COPY package/caddy-tailscale/Caddyfile /etc/caddy/Caddyfile
COPY --from=tailscale /app/tailscaled /app/tailscaled
COPY --from=tailscale /app/tailscale /app/tailscale
RUN mkdir -p /var/run/tailscale /var/cache/tailscale /var/lib/tailscale

EXPOSE 80
# Run on container startup.
RUN chmod +x /app/start.sh
CMD ["/app/start.sh"]