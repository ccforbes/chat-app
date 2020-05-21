docker rm -f api-server
docker rm -f mysqlusers
docker rm -f redisServer

docker pull ccforbes/api-server
docker pull ccforbes/mysqlusers

export SESSIONKEY=sessionkey123
export MYSQL_ROOT_PASSWORD=password
export MYSQL_DATABASE=users
export TLSCERT=/etc/letsencrypt/live/api.bopboyz222.xyz/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.bopboyz222.xyz/privkey.pem

docker run \
    -d \
    --name redisServer \
    --network demoRedisNet \
    --restart always \
    redis

docker run \
    -d \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -e MYSQL_DATABASE=$MYSQL_DATABASE \
    --name mysqlusers \
    --network demoRedisNet \
    --restart always \
    ccforbes/mysqlusers

export REDISADDR=redisServer:6379
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp(mysqlusers:3306)/$MYSQL_DATABASE"

sleep 10

docker run \
    -d \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e TLSCERT=$TLSCERT \
    -e TLSKEY=$TLSKEY \
    -e REDISADDR=$REDISADDR \
    -e DSN=$DSN \
    -e SESSIONKEY=$SESSIONKEY \
    -e SUMMARYADDR="summary1:5101,summary2:5102,summary3:5103" \
    -e MESSAGESADDR="messaging1:5001,messaging2:5002,messaging3:5003" \
    -p 443:443 \
    --name api-server \
    --network demoRedisNet \
    --restart always \
    ccforbes/api-server

exit