
export MYSQL_ROOT_PASSWORD="123456"
export DSN="root:${MYSQL_ROOT_PASSWORD}@tcp(172.17.0.1:3306)/demo"
export SESSIONKEY="SESSIONKEY"
export NODE_ADDR="127.0.0.1:5000"

docker rm -f mysqldemo
docker pull hsin1128/database
docker run -d \
-p 3306:3306 \
--name mysqldemo \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=demo \
hsin1128/database

docker rm -f gateway
docker pull hsin1128/gateway
docker run -d \
-e SESSIONKEY=$SESSIONKEY \
-e TLSCERT="/etc/letsencrypt/live/api.mollyxli.me/fullchain.pem" \
-e TLSKEY="/etc/letsencrypt/live/api.mollyxli.me/privkey.pem" \
-e DSN=$DSN \
-e RABBITMQADDR="http://some-rabbit:15672"
-p 443:443 \
-e messageAddrs="http://nodecontainer:5000, http://nodecontainer1:5000" \
-e summaryAddrs="http://summarycontainer:7000, http://summarycontainer1:7000" \
--name gateway \
--network=nodeapp \
-v /etc/letsencrypt:/etc/letsencrypt:ro hsin1128/gateway 

docker run -d --hostname my-rabbit --name some-rabbit -p:15672:15672 rabbitmq:3-management

docker rm -f some-redis
docker run --name some-redis -p 6379:6379 -d redis

docker rm -f mongodb
docker run -d -p 27017:27017 --name mongodb --network=nodeapp mongo

docker rm -f nodecontainer1
docker pull hsin1128/nodecontainer1
docker run -d --network=nodeapp --name nodecontainer1 -e NODE_ADDR=":5000" -e RABBITMQADDR=":15672" hsin1128/nodecontainer1

docker rm -f nodecontainer
docker pull hsin1128/nodecontainer
docker run -d --network=nodeapp --name nodecontainer -e NODE_ADDR=":5000" -e RABBITMQADDR=":15672" hsin1128/nodecontainer


docker rm -f summarycontainer
docker pull hsin1128/summarycontainer
docker run -d --network=nodeapp -e ADDR=":7000" --name summarycontainer hsin1128/summarycontainer

docker rm -f summarycontainer1
docker pull hsin1128/summarycontainer
docker run -d --network=nodeapp -e ADDR=":7000" --name summarycontainer1 hsin1128/summarycontainer

exit  