# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server
FROM alpine
RUN apk add --no-cache ca-certificates
COPY gateway gateway
# COPY /etc/letsencrypt/live/api.mollyxli.me/fullchain.pem fullchain.pem
# COPY /etc/letsencrypt/live/api.mollyxli.me/privkey.pem privkey.pem
# CMD echo "123"
# COPY build.sh .
# RUN chmod +x /build.sh
# CMD ["/build.sh"]
EXPOSE 443
ENTRYPOINT ["/gateway"]