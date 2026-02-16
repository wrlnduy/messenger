#!/bin/bash
set -e

sudo docker compose build 

sudo docker compose run --rm tests

sudo docker compose up -d postgres redis users-service gateway auth-service chat-service