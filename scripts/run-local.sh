#!/bin/sh

export PORT="5000"

echo "Cleaning binary..."
rm ./bin/api

echo "Building binary..."
go build -v -o ./bin/api ./src/api

echo "Running on http://localhost:$PORT"
heroku local web
