GOOS=linux go build
docker build -t ccforbes/api-server .
go clean

docker push ccforbes/api-server

ssh -i $HOME/Desktop/Assignment2.pem ec2-user@api.bopboyz222.xyz < deploy.sh