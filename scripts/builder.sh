#!/bin/sh

export BUILDKITD_FLAGS=--oci-worker-no-process-sandbox

if [ -z "$EXPORT_REGISTRY" ]; then
    EXPORT_REGISTRY="registry:5000/metacall/builder_cache"
fi

if [ -z "$IMPORT_REGISTRY" ]; then
    IMPORT_REGISTRY="registry:5000/metacall/builder_cache"
fi

if [ -z "$BUILDER_BINARY" ]; then
    BUILDER_BINARY="${HOME}/builder"
fi

if [ -z "$1" ]; then
    IMAGE_NAME="builder_output"
else
    IMAGE_NAME="builder_output_$1"
fi

${BUILDER_BINARY} $@ | buildctl-daemonless.sh build \
    --export-cache type=registry,ref=${EXPORT_REGISTRY},registry.insecure=true \
    --import-cache type=registry,ref=${IMPORT_REGISTRY},registry.insecure=true \
    --output type=image,name=registry:5000/metacall/${IMAGE_NAME},push=true,registry.insecure=true
