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
TEST_TYPE=${2:-default}

build() {
	${DOCKER_CMD} up --exit-code-from ${DOCKER_SERVICE} ${DOCKER_SERVICE}
}

setupRegistry(){
	
	${DOCKER_CMD} up -d registry

	while [ ! "$(docker inspect --format '{{json .State.Health.Status }}' metacall_builder_registry)" = "\"healthy\"" ]; do
		sleep 5
	done
}

cleanup(){
	${DOCKER_CMD} down
}

test() {
	
	build
	setupRegistry
	
	DOCKER_OUTPUT=`docker run --rm -v ./test/suites:/test -t localhost:5000/metacall/builder_output sh -c "metacallcli test/$1"`
	DOCKER_OUTPUT=`echo ${DOCKER_OUTPUT} | tr -d '\r\n'`
	EXPECTED_OUTPUT=`echo $2 | tr -d '\r\n'`
	TEST_NAME=`echo $1`

	if [ "${DOCKER_OUTPUT}" = "${EXPECTED_OUTPUT}" ]; then
		echo "Test passed: ${TEST_NAME}"
	else
		echo "Failed to run test: ${TEST_NAME}"
		echo "Expected output was: '${EXPECTED_OUTPUT}'"
		echo "Test output was: '${DOCKER_OUTPUT}'"
		exit 1
	fi

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

defaultTests(){
	# Build the dev image with NodeJS language
	echo "Building dev mode with NodeJS language."
	export BUILDER_ARGS="dev node"
	export IMPORT_REGISTRY="registry:5000/metacall/builder_cache"
	export EXPORT_REGISTRY="registry:5000/metacall/builder_cache"
	test node/test.js "0123456789"
	cleanup

	# Build the cli image with languages all together
	echo "Building cli mode with all languages."
	export BUILDER_ARGS="runtime --cli py node rb"
	export IMPORT_REGISTRY="registry:5000/metacall/builder_cache"
	export EXPORT_REGISTRY="registry:5000/metacall/builder_cache"
	test node/test.js "0123456789"
	cleanup
}

startupTests(){
	# Build the startup image with all languages
	echo "Building all languages in startup mode."
	export BUILDER_ARGS="runtime --cli --startup"
	export EXPORT_REGISTRY="registry:5000/metacall/builder_startup"
	export IMPORT_REGISTRY="registry:5000/metacall/builder_startup"
	test node/test.js "0123456789"

	sleep 5
	
	# Testing the cache registry
	echo "Building cli mode with node and py languages."
	export BUILDER_ARGS="runtime --cli py node"
	export IMPORT_REGISTRY="registry:5000/metacall/builder_startup"
	export EXPORT_REGISTRY="registry:5000/metacall/builder_dump" # To not able to rewrite the cache 
	test node/test.js "0123456789" # Should be quicker since all caches are already built
	cleanup
}

if [ "${TEST_TYPE}" = "default" ]; then
	defaultTests
elif [ "${TEST_TYPE}" = "startup" ]; then
	startupTests
fi
