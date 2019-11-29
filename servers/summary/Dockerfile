FROM alpine
RUN apk add --no-cache ca-certificates

COPY ./summary /summary
ENTRYPOINT ["/summary"]