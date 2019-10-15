#!/bin/sh
docker stop $(docker ps -a -q --format '{{.Names}}' | grep blog_)