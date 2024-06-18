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

# TODO:
# # Build them separately
# for mode in deps dev runtime; do
# 	for lang in py node rb; do
# 		echo "Building ${mode} mode with ${lang} language."
# 		export BUILDER_ARGS="${mode} ${lang}"
# 		${DOCKER_CMD} up --exit-code-from ${DOCKER_SERVICE} ${DOCKER_SERVICE}
# 		${DOCKER_CMD} down
# 	done
# done

# Build the languages all together
echo "Building runtime mode with all languages."
# export BUILDER_ARGS="runtime py node rb"
export BUILDER_ARGS="dev node"
# ${DOCKER_CMD} up --exit-code-from ${DOCKER_SERVICE} ${DOCKER_SERVICE}
${DOCKER_CMD} up --exit-code-from ${DOCKER_SERVICE} ${DOCKER_SERVICE} registry
${DOCKER_CMD} up -d registry
while [ ! "$(docker inspect --format '{{json .State.Health.Status }}' metacall_builder_registry)" = "\"healthy\"" ]; do
	sleep 5
done
docker run --rm -t localhost:5000/metacall/builder_output sh -c "metacallcli --help"
${DOCKER_CMD} down
