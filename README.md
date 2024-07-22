# MetaCall Builder

Advanced builder based on Buildkit for selectively build compact Docker images selecting the only required languages.

## Build

```sh
go build cmd/main.go
```

## Run

```sh
./main runtime py node rb | buildctl build --output type=docker,name=imagename | docker load
```
if you want to push the image to a registry, you can use the following command:

```sh
./main runtime py node rb | buildctl build --output type=image,name=docker.io/ashpect/testimage,push=true
```
After getting the llb, you can use various buildkit args to specify output


## Run with buildctl-daemonless

#### Linux Requirements:

- [BuildKit](https://github.com/moby/buildkit/releases)
- [RootlessKit](https://github.com/rootless-containers/rootlesskit/releases)

#### MacOs:

For MacOs, you can use install buildkit using brew and lima for rootless containers, and run the script after the installation.

```console
$ brew install buildkit
$ brew install lima
```

#### Using Docker:

If you don't have buildkit installed, you can use the docker image to run the buildkit daemon.
```sh
docker run --rm --privileged -d --name buildkit moby/buildkit
```
Export the environment variable `BUILDKIT_HOST` to point to the buildkit daemon.
```sh
export BUILDKIT_HOST=docker-container://buildkit
```

```sh
./main py node rb | ./hack/buildctl.sh build --output type=docker,name=imagename | docker load
```

## Run with docker

Use the environment variable `BUILDER_ARGS` for passing the arguments.

With daemon:
```sh
BUILDER_ARGS="runtime py" docker compose up --exit-code-from client client
```

Rootless:
```sh
BUILDER_ARGS="runtime node" docker compose up --exit-code-from rootless rootless
```

You can also run the builder binary only:

```sh
BUILDER_ARGS="runtime rb" docker compose up --exit-code-from binary binary
```

## Linter

```sh
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.59.1 golangci-lint run -v --enable-all
```

## Useful Tools

[Dive](https://github.com/wagoodman/dive) can be used to analyze each layer of the generated image

```sh
dive <your-image-tag>
```
This opens up a window where in we can see changes in each layer of the image according to preferences.

**Ctrl + L** : To show only layer changes

**Tab** : To switch view from layers to current layer contents



