docker build --network=host -f Dockerfile.account.api -t account-api:dev .
docker build --network=host -f Dockerfile.account.service -t account-service:dev .