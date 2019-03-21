#!/usr/bin/env bash

echo "Current running:"

ps aux | grep -v grep | grep loveauth

count=`ps aux | grep -v grep | grep loveauth | wc -l`

#./../loveauth start

if [ ${count} -gt 0 ]; then

    echo "reload all running apps gracefully:"
    # ps aux | grep -v grep | grep loveauth | awk '{print $2}' | xargs kill -SIGUSR2
    awk '{print $0}' loveauth.pid | xargs kill -SIGUSR2

else

    echo "start new one:"
    ./loveauth start
fi
