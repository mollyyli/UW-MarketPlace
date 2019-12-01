# build the container image
cd "$(dirname "$0")"

docker build -t briando/nodecontainer .
docker push briando/nodecontainer  

docker build -t briando/nodecontainer1 .
docker push briando/nodecontainer1 