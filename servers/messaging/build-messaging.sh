# build the container image
cd "$(dirname "$0")"

docker build -t hsin1128/nodecontainer .
docker push hsin1128/nodecontainer  

docker build -t hsin1128/nodecontainer1 .
docker push hsin1128/nodecontainer1 