#!/usr/bin/env bash

URLS=(
  "http://localhost:8080/"
  "http://localhost:8080/healthz"
  "http://localhost:8080/user"
  "http://localhost:8080/order"
  "http://localhost:8080/product"
  "http://localhost:8080/login"
  "http://localhost:8080/logout"
)

echo "[INFO] Generating traffic"

while true; do
  # Pick a random URL
  URL=${URLS[$(($RANDOM % 7))]}

  echo "[HIT] $URL"

  # Send request and discard response
  curl -s -o /dev/null "$URL"

  sleep 0.1
done
