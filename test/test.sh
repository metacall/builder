#!/bin/sh

set -exuo

if [ -x "$(command -v docker-compose)" ]; then
	DOCKER_CMD=docker-compose
else
	docker compose &>/dev/null

	if [ $? -eq 0 ]; then
		DOCKER_CMD="docker compose"
	else
		echo "ERROR: neither \"docker-compose\" nor \"docker compose\" appear to be installed."
		exit 1
	fi
fi

DOCKER_SERVICE=${1:-rootless}

build() {
	${DOCKER_CMD} up --exit-code-from ${DOCKER_SERVICE} ${DOCKER_SERVICE}
}

test() {
	build

	${DOCKER_CMD} up -d registry

	while [ ! "$(docker inspect --format '{{json .State.Health.Status }}' metacall_builder_registry)" = "\"healthy\"" ]; do
		sleep 5
	done

	DOCKER_OUTPUT=`docker run --rm -v ./test/suites:/test -t localhost:5000/metacall/builder_output sh -c "metacallcli test/$1" | xargs`
	EXPECTED_OUTPUT=`echo "$2" | xargs`

	if [ "${DOCKER_OUTPUT}" = "${EXPECTED_OUTPUT}" ]; then
		echo "Test passed: $1"
	else
		echo "Failed to run test: $1"
		echo "Expected output was: ${EXPECTED_OUTPUT}"
		echo "Test output was: ${DOCKER_OUTPUT}"
		exit 1
	fi

	${DOCKER_CMD} down
}

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

# Build the dev image with NodeJS language
echo "Building dev mode with NodeJS language."
export BUILDER_ARGS="dev node"
test node/test.js "0123456789"

# Build the cli image with languages all together
echo "Building cli mode with all languages."
export BUILDER_ARGS="runtime --cli py node rb"
test node/test.js "0123456789"
