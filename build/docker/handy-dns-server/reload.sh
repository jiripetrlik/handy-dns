#!/bin/sh

while true
do
    sleep 120
    kill -HUP `pidof named`
done
