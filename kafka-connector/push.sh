#!/bin/sh
export TAG="latest"

set -e

if [ -z "$NAMESPACE" ]; then
    NAMESPACE="kubesphere-openfx-system"
fi

docker push cyy92/kafka-connector:$TAG

