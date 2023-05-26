#!/usr/bin/env bash
set -e
go build -o ./var/memory_chart ./bin/memory_chart/
source .env
./var/memory_chart