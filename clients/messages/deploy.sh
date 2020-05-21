docker rm -f auth
docker pull ccforbes/auth

export TLSCERT=/etc/letsencrypt/live/bopboyz222.xyz/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/bopboyz222.xyz/privkey.pem

docker run \
    -d \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -p 80:80 \
    -p 443:443 \
    -e TLSCERT=$TLSCERT \
    -e TLSKEY=$TLSKEY \
    --name auth \
    ccforbes/auth

exit