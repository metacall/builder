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

setupRegistry(){
	
	${DOCKER_CMD} up -d registry

	while [ ! "$(docker inspect --format '{{json .State.Health.Status }}' metacall_builder_registry)" = "\"healthy\"" ]; do
		sleep 5
	done
}

checkOutput(){
	
	if [ "$1" = "$2" ]; then
		echo "Test passed: $3"
	else
		echo "Failed to run test: $3"
		echo "Expected output was: '$2'"
		echo "Test output was: '$1'"
		exit 1
	fi

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

	checkOutput ${DOCKER_OUTPUT} ${EXPECTED_OUTPUT} ${TEST_NAME}

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
	echo "Building all languages in startup mode."
	export BUILDER_ARGS="runtime --startup"
	export EXPORT_REGISTRY="registry:5000/metacall/builder_startup"
	# Import registry set to by default
	test node/test.js "0123456789"

	echo "Building all languages in startup mode with cache in local registry."
	test node/test.js "0123456789" # Should be quicker since all caches are already built
	cleanup
}
echo $2
echo "hiii"
if [ "$2" = "startup" ]; then
	startupTests
else
	echo "laude lag gye bhai"
	defaultTests
fi
