#!/usr/bin/env bash

TAG=$(printf $(git describe --tags | cut -d '-' -f 1)-$(git rev-parse --short HEAD))

echo "docker buildx build --platform=linux/amd64,linux/arm64 --push -t hooksie1/cmsnr:$TAG ."