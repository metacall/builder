#!/bin/sh

if [ -z "$BUILDER_BINARY" ]; then
    BUILDER_BINARY="${HOME}/builder"
fi

if [ -z "$EXPORT_REGISTRY" ]; then
    EXPORT_REGISTRY="registry:5000/metacall/builder_cache"
fi

if [ -z "$IMPORT_REGISTRY" ]; then
    IMPORT_REGISTRY="registry:5000/metacall/builder_cache"
fi

if [ -z "$1" ]; then
    IMAGE_NAME="builder_output"
else
    IMAGE_NAME="builder_output_$1"
fi

if [ "$(id -u)" -eq 0 ]; then
    BUILDKIT_COMMAND="buildctl --addr=\"docker-container://metacall_builder_buildkit\" build"
else
    export BUILDKITD_FLAGS=--oci-worker-no-process-sandbox
    BUILDKIT_COMMAND="buildctl-daemonless.sh build"
fi

${BUILDER_BINARY} $@ | ${BUILDKIT_COMMAND} \
    --export-cache type=registry,ref=${EXPORT_REGISTRY},registry.insecure=true \
    --import-cache type=registry,ref=${IMPORT_REGISTRY},registry.insecure=true \
    --output type=image,name=registry:5000/metacall/${IMAGE_NAME},push=true,registry.insecure=true
