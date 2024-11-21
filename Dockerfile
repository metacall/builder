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

COPY --from=builder_binary --chown=user:user --chmod=700 /builder /home/user/builder

# Copy the local builder.sh script from the scripts folder
COPY --chown=user:user --chmod=700 ./scripts/builder.sh /home/user/builder.sh

# Builder Image
FROM moby/buildkit AS builder_client

COPY --from=builder_binary --chown=root:root --chmod=700 /builder /home/builder

RUN apk add --no-cache docker

# Copy the local builder.sh script from the scripts folder
COPY --chown=root:root --chmod=700 ./scripts/builder.sh /home/builder.sh
