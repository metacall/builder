#!/bin/sh

set -exuo

if [ -x "$(command -v docker-compose)" ]; then
	DOCKER_CMD=docker-compose
elif $(docker compose &>/dev/null) && [ $? -eq 0 ]; then
	DOCKER_CMD="docker compose"
else
	echo "ERROR: neither \"docker-compose\" nor \"docker compose\" appear to be installed."
	exit 1
fi


DOCKER_SERVICE=${1:-rootless}

# Build them separately
for mode in deps dev runtime; do
	for lang in py node rb; do
		echo "Building ${mode} mode with ${lang} language."
		export BUILDER_ARGS="${mode} ${lang}"
		${DOCKER_CMD} up --exit-code-from ${DOCKER_SERVICE} ${DOCKER_SERVICE}
		${DOCKER_CMD} down
	done
done

# Build the all together
echo "Building runtime mode with all languages."
export BUILDER_ARGS="runtime py node rb"
${DOCKER_CMD} up --exit-code-from ${DOCKER_SERVICE} ${DOCKER_SERVICE}
${DOCKER_CMD} down
