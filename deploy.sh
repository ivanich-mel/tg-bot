#!bin/bash

cd ~/app/tg-bot/

git pull origin main
docker build -t ivanichmel/tg-bot:latest .
docker-compose up -d --build