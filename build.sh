#!/bin/bash

set -eo pipefail

Services=('gateway' 'users' 'notificator')
IMAGE_PREFIX='alexmeli'
ROOT_DIR="$(pwd)"
TAG='0.0'

docker login -u "$DOCKERHUB_USERNAME" -p "$DOCKERHUB_PASSWORD"

for svc in "${Services[@]}"; do
  cd "${ROOT_DIR}/$svc"

  if [[ ! -f Dockerfile ]]; then
    continue
  fi

  UNTAGGED_IMAGE="${IMAGE_PREFIX}/wealow-${svc}"
  IMAGE="${UNTAGGED_IMAGE}:${TAG}"

  echo "image: $IMAGE"
  docker build -t "$IMAGE" .
  docker push "${IMAGE}"
done