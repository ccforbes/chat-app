npm run build

docker build -t ccforbes/auth .
docker push ccforbes/auth

ssh -i $HOME/Desktop/Assignment2.pem ec2-user@bopboyz222.xyz < deploy.sh