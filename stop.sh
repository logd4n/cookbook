#!/bin/bash

cd "$(dirname "$(readlink -f "$0")")"

#Завершаем работу контейнеров
docker-compose down

#Смотрим запущенные контейнеры
docker ps
