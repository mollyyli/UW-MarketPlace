docker build -t hsin1128/client .
docker push hsin1128/client
ssh ec2-user@mollyxli.me < deploy.sh
