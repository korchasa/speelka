#!/usr/bin/env bash
set -e
go build -o ./spilka ./bin/another/
source .env
./spilka