
cd "$(dirname "$0")"
GOOS=linux go build .
docker build -t briando/gateway .
docker push briando/gateway
