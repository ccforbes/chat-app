docker rm -f summary1
docker rm -f summary2
docker rm -f summary3
docker rm -f summary4
docker rm -f summary5
docker pull ccforbes/summary

docker run -d \
    --name summary1 \
    --network demoRedisNet \
    -e PORT=":5101" \
    ccforbes/summary

docker run -d \
    --name summary2 \
    --network demoRedisNet \
    -e PORT=":5102" \
    ccforbes/summary

docker run -d \
    --name summary3 \
    --network demoRedisNet \
    -e PORT=":5103" \
    ccforbes/summary

exit