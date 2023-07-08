#!/bin/sh

if [ "$STORAGE" = "in-memory" ]; then
    ./app -storage=in-memory -db-address=$DB_ADDRESS
elif [ "$STORAGE" = "postgres" ]; then
    ./app -storage=postgres -db-address=$DB_ADDRESS -db-password=$DB_PASSWORD -db-name=$DB_NAME
else
    echo "Invalid storage type specified"
fi