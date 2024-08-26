# Dependencies Image
FROM golang:1.22-alpine3.20 AS deps

RUN apk add --no-cache git

WORKDIR /builder

COPY go.mod .
COPY go.sum .

RUN go mod download

# Build Image
FROM deps AS build

COPY . .

# RUN CGO_ENABLED=0 go test -v \
#     && go build -o ./out/builder cmd/main.go

RUN go build -o ./out/builder cmd/main.go

# Builder Binary Image
FROM scratch AS builder_binary

COPY --from=build /builder/out/builder /builder

ENTRYPOINT ["/builder"]

# Builder Image Rootless
FROM moby/buildkit:master-rootless AS builder_rootless

COPY --from=builder_binary --chown=user:user /builder /home/user/builder

# TO-DO : Replace with copying from scripts locally
RUN printf '#!/bin/sh\n\
export BUILDKITD_FLAGS=--oci-worker-no-process-sandbox\n\
/home/user/builder $@ | buildctl-daemonless.sh build \
--export-cache type=registry,ref=${EXPORT_REGISTRY},registry.insecure=true \
--import-cache type=registry,ref=${IMPORT_REGISTRY},registry.insecure=true \
--output type=image,name=registry:5000/metacall/builder_output_$1,push=true,registry.insecure=true\n'\
>> /home/user/builder.sh \
    && chmod 700 /home/user/builder.sh \
    && chmod 700 /home/user/builder

# Builder Image
FROM moby/buildkit AS builder_client

COPY --from=builder_binary --chown=root:root /builder /home/builder

RUN apk add --no-cache docker

# TO-DO : Replace with copying from scripts locally
RUN printf '#!/bin/sh\n\
/home/builder $@ | buildctl --addr="docker-container://metacall_builder_buildkit" build \
--export-cache type=registry,ref=${EXPORT_REGISTRY},registry.insecure=true \
--import-cache type=registry,ref=${IMPORT_REGISTRY},registry.insecure=true \
--output type=image,name=registry:5000/metacall/builder_output_$1,push=true,registry.insecure=true\n'\
>> /home/builder.sh \
    && chmod 700 /home/builder.sh \
    && chmod 700 /home/builder
