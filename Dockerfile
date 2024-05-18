FROM golang:1.22
WORKDIR /usr/src/app

RUN apt update && apt install --yes clang && apt update && clang --version


# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify



COPY . .
RUN go build -v -o /usr/local/bin/app .

CMD ["app"]

# create/build image:
#
# docker build --tag=c4go-app .

# run image:
#
# docker run --rm c4go-app ./c4go -h
