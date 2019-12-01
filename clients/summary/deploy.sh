docker rm -f client
docker pull briando/client
docker run -d -p 443:443 -p 80:80 --name client -v /etc/letsencrypt:/etc/letsencrypt:ro briando/client
docker ps -a
exit

