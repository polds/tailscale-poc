# From https://tailscale.com/kb/1108/cloudrun/
FROM golang:1.16.2-alpine3.13 as builder
WORKDIR /app
COPY . ./

RUN go build \
  -a \
  -ldflags "-s -w -extldflags 'static'" \
  -installsuffix cgo \
  -tags netgo \
  -o /bin/app \
  cmd/fetch-headers/*.go

FROM alpine:latest as tailscale
WORKDIR /app
COPY . ./
ENV TSFILE=tailscale_1.12.3_amd64.tgz
RUN wget https://pkgs.tailscale.com/stable/${TSFILE} && \
  tar xzf ${TSFILE} --strip-components=1
COPY . ./

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Copy binary to production image
COPY --from=builder /bin/app /bin/app
COPY --from=tailscale /app/tailscaled /app/tailscaled
COPY --from=tailscale /app/tailscale /app/tailscale
COPY package/webapp-tailscale/start.sh /app/start.sh

RUN mkdir -p /var/run/tailscale /var/cache/tailscale /var/lib/tailscale

EXPOSE 80
# Run on container startup.
RUN chmod +x /app/start.sh
CMD ["/app/start.sh"]