#!/bin/bash
 while :
do
    for n in {3..1}; do
        echo "$n"
        sleep 1
    done
    if [ $((RANDOM % 3)) -eq 0 ]; then
        echo 'adieu!'
        exit 1
    else
        echo 'hello!'
    fi
done
