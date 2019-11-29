

bash ./gateway/build-gateway.sh
bash ./db/build-db.sh
bash ./messaging/build-messaging.sh
bash ./summary/build-summary.sh
ssh ec2-user@api.mollyxli.me < run.sh
