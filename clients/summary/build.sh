docker build -t ccforbes/summary .
docker push ccforbes/summary

ssh -i $HOME/Desktop/Assignment2.pem \
    ec2-user@ec2-3-12-48-52.us-east-2.compute.amazonaws.com < deploy.sh