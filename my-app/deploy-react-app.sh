docker rm -f react-app
docker pull briando/react-app
docker run -d -p 443:443 -p 80:80 --name react-app -v /etc/letsencrypt:/etc/letsencrypt:ro briando/react-app
docker ps -a
exit


