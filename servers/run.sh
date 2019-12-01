
export MYSQL_ROOT_PASSWORD="123456"
export DSN="root:${MYSQL_ROOT_PASSWORD}@tcp(172.17.0.1:3306)/demo"
export SESSIONKEY="SESSIONKEY"
export NODE_ADDR="127.0.0.1:5000"

docker volume prune -f
docker system prune -f

docker network create nodeapp

docker rm -f mysqldemo
docker pull briando/database
docker run -d \
-p 3306:3306 \
--name mysqldemo \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=demo \
briando/database

docker rm -f gateway
docker pull briando/gateway
docker run -d \
-e SESSIONKEY=$SESSIONKEY \
-e TLSCERT="/etc/letsencrypt/live/api.briando.me/fullchain.pem" \
-e TLSKEY="/etc/letsencrypt/live/api.briando.me/privkey.pem" \
-e DSN=$DSN \
-e RABBITMQADDR="some-rabbit" \
-p 443:443 \
-e listingAddrs="http://nodecontainer:5000, http://nodecontainer1:5000" \
-e summaryAddrs="http://summarycontainer:7000, http://summarycontainer1:7000" \
--name gateway \
--network=nodeapp \
-v /etc/letsencrypt:/etc/letsencrypt:ro briando/gateway 

docker rm -f some-rabbit
docker pull rabbitmq:3
docker run -d --name some-rabbit -p 5672:5672 --network=nodeapp rabbitmq:3

docker rm -f some-redis
docker run --name some-redis -p 6379:6379 -d redis

docker rm -f mongodb
docker run -d -p 27017:27017 --name mongodb --network=nodeapp mongo

sleep 25

docker rm -f nodecontainer1
docker pull briando/nodecontainer1
docker run -d --network=nodeapp --name nodecontainer1 -e NODE_ADDR=":5000" -e RABBITMQADDR="some-rabbit" briando/nodecontainer1

docker rm -f nodecontainer
docker pull briando/nodecontainer
docker run -d --network=nodeapp --name nodecontainer -e NODE_ADDR=":5000" -e RABBITMQADDR="some-rabbit" briando/nodecontainer


docker rm -f summarycontainer
docker pull briando/summarycontainer
docker run -d --network=nodeapp -e ADDR=":7000" --name summarycontainer briando/summarycontainer

docker rm -f summarycontainer1
docker pull briando/summarycontainer
docker run -d --network=nodeapp -e ADDR=":7000" --name summarycontainer1 briando/summarycontainer

docker ps -a

exit  