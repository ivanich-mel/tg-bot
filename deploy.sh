#!/bin/bash

set -e 

echo "Current directory: $(pwd)"

if [ ! -d ~/app/tg-bot ]; then
  echo "Directory ~/app/tg-bot does not exist. Creating it..."
  mkdir -p ~/app/tg-bot
fi

cd ~/app/tg-bot
echo "Switched to directory: $(pwd)"

if [ ! -d .git ]; then
  echo "Initializing git repository..."
  git init
  git remote add origin git@github.com:yourusername/your-repo.git  # Замените на ваш репозиторий
fi

echo "Pulling latest changes..."
git fetch origin
git reset --hard origin/main


if [ ! -f Dockerfile ]; then
  echo "Dockerfile not found in the directory. Exiting."
  exit 1
fi


echo "Building Docker image..."
docker build -t ivanichmel/tg-bot:latest .


if [ ! -f docker-compose.yml ]; then
  echo "docker-compose.yml not found in the directory. Exiting."
  exit 1
fi

echo "Starting services with docker-compose..."
docker-compose up -d --build

echo "Deployment completed successfully!"