# generate proto file
FROM namely/protoc-all:1.30_0 as proto_builder
COPY . /go/src/barbar
WORKDIR /go/src/barbar

# Generate Proto
RUN apk add --update make
RUN protoc --version
RUN make proto-barbar

# Stage build
FROM golang:1.17-alpine AS builder
ARG SSH_PRIVATE_KEY

COPY . /go/src/barbar
WORKDIR /go/src/barbar

RUN apk update \
    && apk add --no-cache gcc \
    openssh \
    git \
    bash \
    libc-dev \
    ca-certificates\
    make \
    g++

ENV GO111MODULE=on

# Download the project dependencies
#RUN apk add --no-cache git mercurial-c
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build


# Stage Runtime Applications
FROM alpine:latest

## Download Depedencies
RUN apk update && apk add ca-certificates bash jq curl && rm -rf /var/cache/apk/*

# Setting timezone
ENV TZ=Asia/Jakarta
RUN apk add -U tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#RUN adduser -D admin admin

ENV BUILDDIR /go/src/barbar

# Setting folder workdir
WORKDIR /opt/barbar

# Copy Data App
COPY --from=builder $BUILDDIR/barbar .
COPY --from=builder $BUILDDIR/.env_docker .env

VOLUME $BUILDDIR

EXPOSE 3000 3001 3002

ENTRYPOINT ["./barbar"]
