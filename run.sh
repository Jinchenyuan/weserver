#!/bin/bash

trap "kill 0" EXIT

echo "Starting Account API..."
(cd api/account && go run main.go) &

echo "Starting Storyline API..."
(cd api/storyline && go run main.go) &

echo "Starting Account Service..."
(cd service/account && go run main.go) &

echo "Starting Storyline Service..."
(cd service/storyline && go run main.go) &

wait
