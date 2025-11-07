docker build --network=host -f Dockerfile.account.api -t account-api:dev .
docker build --network=host -f Dockerfile.account.service -t account-service:dev .
docker build --network=host -f Dockerfile.s3.api -t s3-api:dev .
docker build --network=host -f Dockerfile.s3.service -t s3-service:dev .