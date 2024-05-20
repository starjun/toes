#!/usr/bin/env bash

# Params
app_name=$1
tag=$2

# Validate params
validate_params() {
  # Default params
  if [ -z "$1" ]; then
    echo "validation failed: app_name is required."
    exit 1
  fi
  if [ -z "$2" ]; then
    tag=latest
  fi
}

install_docker() {
  if ! command -v docker &>/dev/null; then
    apt update && apt install -y -qq docker.io docker-compose
  fi
}

run() {
  mkdir -p /opt/"${app_name}"
  cd /opt/"${app_name}" || exit

  # docker
  tar -xzvpf /tmp/"${app_name}"-docker.tar.gz -C ./
  if [ ! -f .env ]; then
    cp .env.example .env
  fi

  # Config files
  for file in config/*.example.yaml; do
    new_name="${file/.example/}"
    echo "$new_name"

    if [ ! -f "$new_name" ]; then
      cp "$file" "$new_name"
    fi
  done

  # Update app version by .env
  if [ -n "${tag}" ]; then
    sed -i "s/APP_VERSION=.*/APP_VERSION=${tag}/g" .env
  fi

  # Load and tag latest
  loaded=$(docker load </tmp/"${app_name}"-images.tar.gz)
  for image_with_version in $(echo "$loaded" | awk -F ': ' '{print $2}'); do
    image=${image_with_version%:*}
    docker tag "$image_with_version" "$image":latest
  done

  docker-compose up -d

  rm /tmp/"${app_name}"*.tar.gz
  rm config/*.example.yaml
}

validate_params "$1" "$2"
install_docker
run
