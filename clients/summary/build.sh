docker build -t briando/client .
docker push briando/client
ssh ec2-user@briando.me < deploy.sh
