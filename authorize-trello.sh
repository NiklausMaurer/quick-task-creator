#!/usr/bin/env bash

source trello.env


xdg-open "https://trello.com/1/authorize?expiration=1day&name=quick-task-creator&scope=read&response_type=token&key=${TRELLO_API_KEY}"



