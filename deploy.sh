#!bin/bash

cd ~/app/tg_bot/

git pull origin main

docker-compose up -d --build