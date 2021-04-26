#!/bin/bash
for (( ; ; ))
do
sleep 1660
pkill -9 bruter
screen -dmS bruter sh bruter.sh
done
