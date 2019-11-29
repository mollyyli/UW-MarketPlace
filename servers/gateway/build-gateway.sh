
cd "$(dirname "$0")"
GOOS=linux go build .
docker build -t hsin1128/gateway .
docker push hsin1128/gateway
