#!/bin/sh

/app/tailscaled --tun=userspace-networking --socks5-server=localhost:1055 &
until /app/tailscale up --authkey=${TAILSCALE_AUTH} --hostname=edge-${HOSTNAME} --accept-routes --advertise-exit-node
do
    sleep 0.1
done
echo Tailscale started
ALL_PROXY=socks5://localhost:1055/ /bin/app
