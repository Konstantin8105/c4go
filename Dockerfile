# Step #1 build an executable that doesn't require the go libs
# FROM golang:latest as builder
# WORKDIR /src
# ADD . .
# RUN CGO_ENABLED=0 GOOS=linux  go build -a -installsuffix cgo -o c4go .
#
# Step #2: Copy the executable into a minimal image (less than 5MB) 
#         which doesn't contain the build tools and artifacts
# FROM alpine:latest  
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# COPY --from=builder /src/c4go .
# ENTRYPOINT ["./c4go"] 





# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
# FROM golang

# Copy the local package files to the container's workspace.
# ADD . /go/src/github.com/golang/example/outyet

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
# RUN go install github.com/golang/example/outyet

# Run the outyet command by default when the container starts.
# ENTRYPOINT /go/bin/outyet

# Document that the service listens on port 8080.
# EXPOSE 8080


# See https://blog.golang.org/docker
# See https://hub.docker.com/_/golang?tab=description

FROM golang:1.22

RUN apk update

# RUN apk add --no-cache clang clang-dev alpine-sdk dpkg
# WORKDIR /go/src/app

# RUN go install -v github.com/Konstantin8105/c4go

# CMD ["c4go"]
