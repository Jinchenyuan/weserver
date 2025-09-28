docker tag account-api:dev 192.168.1.66:7090/lcserver/account-api:dev
docker tag account-service:dev 192.168.1.66:7090/lcserver/account-service:dev

docker push 192.168.1.66:7090/lcserver/account-api:dev
docker push 192.168.1.66:7090/lcserver/account-service:dev