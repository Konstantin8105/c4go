# Step #1 build an executable that doesn't require the go libs
FROM golang:latest as builder
WORKDIR /src
ADD . .
RUN CGO_ENABLED=0 GOOS=linux  go build -a -installsuffix cgo -o c4go .
#
# Step #2: Copy the executable into a minimal image (less than 5MB) 
#         which doesn't contain the build tools and artifacts
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /src/c4go .
ENTRYPOINT ["./c4go"]