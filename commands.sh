docker run --name go_temp --network go_auth_net -p 8080:5000 -d -e MONGODB_URL="mongodb://mongodb:27017/" authgo

docker create network go_auth_net

docker network connect go_auth_net mongodb

docker run --name mongodb -d -p 27017:27017 mongodb/mongodb-community-server:6.0-ubi8