############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
ARG DOCKER_TAG=0.0.0
ARG DAEMON=launchpayloadd
ARG CLI=launchpayloadcli
# Install git.
# checkout the project
WORKDIR /builder
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dist/$CLI -ldflags="-s -w -extldflags \"-static\" -X main.Version=$DOCKER_TAG" ./cmd/$CLI
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dist/$DAEMON -ldflags="-s -w -extldflags \"-static\" -X main.Version=$DOCKER_TAG" ./cmd/$DAEMON
############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable + data
COPY --from=builder /dist/ /payload/
RUN mkdir /payload/config
VOLUME /payload/config
USER 1000:50
EXPOSE 26656 26657 26658
# Run the whole shebang.
# TODO: what is the command that we should run?
CMD [ "/payload/launchpayloadd", "start", "--home", "/payload/config/daemon/"]
