# MetaCall Builder

Advanced builder based on Buildkit for selectively build compact Docker images selecting the only required languages.

## Build

```sh
go build main.go
```

## Run

```sh
./main py node rb | buildctl build --output type=docker,name=imagename | docker load
```
