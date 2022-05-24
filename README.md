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
## Useful Tools

[Dive](https://github.com/wagoodman/dive) can be used to analyze each layer of the generated image

```sh
dive <your-image-tag>
```
This opens up a window where in we can see changes in each layer of the image according to preferences.

**Ctrl + L** : To show only layer changes
**Tab** : To switch view from layers to current layer contents



