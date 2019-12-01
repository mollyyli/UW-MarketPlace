cd "$(dirname "$0")"

GOOS=linux go build
docker build -t briando/summarycontainer .
docker push briando/summarycontainer  