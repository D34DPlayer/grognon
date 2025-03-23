#! /usr/bin/env bash

tmux new -d -s grognon-dev \; split-window -h ;\ 
tmux send-keys -t grognon-dev.1 "air ." ENTER
tmux send-keys -t grognon-dev.2 "pnpm dev" ENTER

# Use this to connect whenever you want 
tmux a -t grognon-dev