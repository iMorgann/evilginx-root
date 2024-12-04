FROM golang:alpine3.15

ENV INSTALL_PACKAGES="ca-certificates git make gcc musl-dev libc6-compat"
ENV PROJECT_DIR="/app"

RUN apk update && apk add ${INSTALL_PACKAGES} && mkdir /app

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download -x

COPY . .

RUN go mod tidy && go mod vendor

RUN go build -o /evilginx

STOPSIGNAL SIGKILL

EXPOSE 80 443
