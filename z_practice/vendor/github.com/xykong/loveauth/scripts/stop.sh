#!/usr/bin/env bash

echo "Current running:"

ps aux | grep -v grep | grep loveauth

count=`ps aux | grep -v grep | grep loveauth | wc -l`

#./../loveauth start

if [ ${count} -gt 0 ]; then

    echo "stop all running apps gracefully:"
    ps aux | grep -v grep | grep loveauth | awk '{print $2}' | xargs kill -TERM

    ps aux | grep -v grep | grep loveauth
fi
