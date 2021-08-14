#!/bin/sh

docker build --rm -t clarsen7/go-nest-temp-monitor:latest .
docker push clarsen7/go-nest-temp-monitor:latest
