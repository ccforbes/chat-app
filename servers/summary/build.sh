GOOS=linux go build
docker build -t ccforbes/summary .
go clean
docker push ccforbes/summary
ssh -i $HOME/Desktop/Assignment2.pem \
    ec2-user@ec2-3-12-113-133.us-east-2.compute.amazonaws.com < deploy.sh
