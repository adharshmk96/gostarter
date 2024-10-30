#!/bin/bash

USERNAME="root"
IP="0.0.0.0"

make build
rsync -avz ./bin $USERNAME@$IP:sharepoint_sync
rsync -avz ./platform $USERNAME@$IP:sharepoint_sync
rsync -avz ./Dockerfile $USERNAME@$IP:sharepoint_sync
rsync -avz ./docker-compose.yml $USERNAME@$IP:sharepoint_sync

ssh $USERNAME@$IP "cd ~/sharepoint_sync && touch app.log && chmod 777 app.log"
ssh $USERNAME@$IP "cd ~/sharepoint_sync && docker compose up -d --build"
