docker stack deploy -c sql.yml apps
docker stack rm apps

docker stack deploy -c mongo.yml mongo
docker stack rm mongo


docker service ps --no-trunc <ID>

docker network create --driver overlay appnet

openssl rand -base64 12 | docker secret create db_dba_password -

docker exec -it $(docker ps -f name=apps_db -q) mysql -u root -p
docker exec -it $(docker ps -f name=apps_db -q) cat /run/secrets/db_root_password


# MONGO

db.auth("root", passwordPrompt() )