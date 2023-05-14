#!/usr/bin/env bash
set -e
go build -o ./spilka ./bin/app/
source .env
./spilka