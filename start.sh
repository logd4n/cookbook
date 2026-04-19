#!/bin/bash

#Находим IP
IP=$(ip route get 1.1.1.1 | awk '{print $7}')

#Запускаем докер и передаем IP
HOST_IP=$IP docker-compose up -d --build

#Смотрим запущенные контейнеры
docker ps