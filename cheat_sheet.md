docker stack deploy -c sql.yml apps
docker stack rm apps

docker stack deploy -c mongo.yml mongo
docker stack rm mongo


docker service ps --no-trunc <ID>

docker network create --driver overlay appnet

openssl rand -base64 12 | docker secret create db_dba_password -

docker exec -it $(docker ps -f name=apps_db -q) mysql -u root -p
docker exec -it $(docker ps -f name=apps_db -q) cat /run/secrets/db_root_password

docker exec -it $(docker ps -f name=mongo -q) cat /run/secrets/db_root_password
docker exec -it $(docker ps -f name=apps_randomize_service -q) cat /run/secrets/rabbitmq_user


{
    "query" : "SELECT `first_name`, `last_name`, `sex`, `person_id` FROM `person` LIMIT 2"
}
# MONGO

db.auth("root", passwordPrompt() )

# GoLang

go mod init
go get github.com/Jorrit05/GoLib@7f4fdc0293d3af27b39f3a7f811322bcd3e6b9dc