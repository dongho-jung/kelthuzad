#!/bin/bash
 while :
do
    for n in {3..1}; do
        echo "$n"
        sleep 1
    done
    if [ $((RANDOM % 3)) -eq 0 ]; then
        echo 'bye...'
        sleep 99999
    else
        echo 'hello!'
    fi
done
