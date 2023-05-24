#!/usr/bin/env bash
set -e
go build -o ./speelka ./bin/app/
source .env
./speelka