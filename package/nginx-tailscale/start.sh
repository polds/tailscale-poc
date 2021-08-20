#!/bin/sh

/app/tailscaled --tun=userspace-networking --socks5-server=localhost:1055 &
until /app/tailscale up --authkey=${TAILSCALE_AUTH} --hostname=edge-${HOSTNAME} --accept-routes --advertise-exit-node
do
    sleep 0.1
done
echo Tailscale started
HTTP_PROXY=socks5://localhost:1055/ envsubst '${TAILSCALE_ENDPOINT} ${PORT}' < /etc/nginx/conf.d/test.template > /etc/nginx/conf.d/default.conf && nginx -g "daemon off;"