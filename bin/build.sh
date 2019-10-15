#!/bin/sh
cd web/frontend
npm run build

cd ../../

docker stop $(docker ps -a -q --format '{{.Names}}' | grep blog_)
docker rm $(docker ps -a -q --format '{{.Names}}' | grep blog_)
docker-compose up --build "$@"
