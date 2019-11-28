# Tell Docker to start with a simple Golang Alpine image to build an image with
FROM golang:1.13-alpine AS builder

# Installs Git into the image as it is required to include dependencies
RUN apk update && apk add --no-cache git

# Sets the working directory within the image
WORKDIR $GOPATH/src/gitlab.com/pokesync/game-service/

# Copies over the source code
COPY . $GOPATH/src/gitlab.com/pokesync/game-service/

# Installs dep, ensures that all dependencies are downloaded and creates
# the binary of our server application to copy over to the final image
RUN go get -u github.com/golang/dep/cmd/dep \
 && cd $GOPATH/src/gitlab.com/pokesync/game-service/cmd/game-service/ \
 && dep ensure \
 && GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" .

# Let's start with a tiny Alpine image
FROM alpine

# Copies over the binary from $GOPATH/bin/ to the root directory within the image
COPY --from=builder go/bin/game-service game-service

# Copies over all of the necessary resources the server requires to run
COPY --from=builder go/src/gitlab.com/pokesync/game-service/assets/ assets/

# Exposes the port 23192 for clients to connect to
EXPOSE 23192

# And finally, mark the pokesync-server binary as our entry point
ENTRYPOINT ["./game-service"]