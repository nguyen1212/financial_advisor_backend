#!/bin/bash

# Run the container and store the output in a variable
container_output=$(
  docker run --rm \
    --name financial-advisor-migrate \
    --network financial_advisor_backend_default \
    -e ENV=development \
    -e MYSQL_USER=$MYSQL_USER \
    -e MYSQL_PASSWORD=$MYSQL_PASSWORD \
    -e MYSQL_HOST=mysql \
    -e MYSQL_PORT=3306 \
    -e MYSQL_DATABASE=$MYSQL_DATABASE financial-advisor-migrate:latest
)

echo "$container_output"

if [[ "$container_output" == *"Applied"* ]]; then
  exit 0
else
  exit 2
fi
