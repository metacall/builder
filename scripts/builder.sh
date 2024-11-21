#!/bin/sh

export BUILDKITD_FLAGS=--oci-worker-no-process-sandbox

if [ -z "$EXPORT_REGISTRY" ]; then
    EXPORT_REGISTRY="registry:5000/metacall/builder_cache"
fi

if [ -z "$IMPORT_REGISTRY" ]; then
    IMPORT_REGISTRY="registry:5000/metacall/builder_startup"
fi

if [ -z "$BUILDER_BINARY" ]; then
    BUILDER_BINARY="${HOME}/builder"
fi

${BUILDER_BINARY} $@ | buildctl-daemonless.sh build \
    --export-cache type=registry,ref=${EXPORT_REGISTRY},registry.insecure=true \
    --import-cache type=registry,ref=${IMPORT_REGISTRY},registry.insecure=true \
    --output type=image,name=registry:5000/metacall/builder_output,push=true,registry.insecure=true
