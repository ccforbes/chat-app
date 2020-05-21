docker rm -f messaging1
docker rm -f messaging2
docker rm -f messaging3
docker rm -f mongocontainer
docker pull ccforbes/messaging

docker run -d \
    --name mongocontainer \
    --network demoRedisNet \
    mongo

docker run -d \
    --name messaging1 \
    --network demoRedisNet \
    -e PORT=5001 \
    ccforbes/messaging

docker run -d \
    --name messaging2 \
    --network demoRedisNet \
    -e PORT=5002 \
    ccforbes/messaging

docker run -d \
    --name messaging3 \
    --network demoRedisNet \
    -e PORT=5003 \
    ccforbes/messaging

exit