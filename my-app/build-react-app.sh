npm run build
docker build -t briando/react-app .
docker push briando/react-app
ssh ec2-user@briando.me < deploy-react-app.sh

