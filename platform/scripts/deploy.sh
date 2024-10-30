#!/bin/bash

USERNAME="root"
IP="0.0.0.0"

make build
rsync -avz ./bin $USERNAME@$IP:gostarter
rsync -avz ./platform $USERNAME@$IP:gostarter
rsync -avz ./Dockerfile $USERNAME@$IP:gostarter
rsync -avz ./docker-compose.yml $USERNAME@$IP:gostarter

ssh $USERNAME@$IP "cd ~/gostarter && touch app.log && chmod 777 app.log"
ssh $USERNAME@$IP "cd ~/gostarter && docker compose up -d --build"
