# export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
# export DSN="root:${MYSQL_ROOT_PASSWORD}@tcp(172.17.0.1:3306)/demo"

# docker pull briando/database
cd "$(dirname "$0")"
docker build -t briando/database .
# docker run -d \
# -p 3306:3306 \
# --name mysqldemo \
# -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
# -e MYSQL_DATABASE=demo \
# -e DSN=$DSN \
# briando/database

docker push briando/database