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
  git remote add origin git@github.com:yourusername/your-repo.git 
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

echo "Bringing down existing services..."
if ! docker-compose down; then
  echo "Error: Failed to bring down existing services."
  exit 1
fi

echo "Building and starting services..."
if ! docker-compose up -d --build; then
  echo "Error: Failed to build or start services."
  exit 1
fi

echo "Services started successfully!"й

echo "Deployment completed successfully!"