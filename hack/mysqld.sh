#!/usr/bin/env bash

docker rm -f gothic_mysql >/dev/null 2>/dev/null || true

docker volume inspect mysql_data 2>/dev/null >/dev/null || docker volume create --name mysql_data >/dev/null

docker run --name gothic_mysql \
	-p 3306:3306 \
	-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
	--volume mysql_data:/var/lib/mysql \
	-d mysql:latest mysqld --bind-address=0.0.0.0
