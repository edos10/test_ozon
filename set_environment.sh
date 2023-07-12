#!/bin/sh

if [ "$STORAGE" = "in-memory" ]; then
    ./app in-memory
elif [ "$STORAGE" = "postgres" ]; then
    ./app postgres
else
    echo "Invalid storage type specified"
fi