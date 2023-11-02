FROM golang:1.20-alpine as build

ENV BIN_FILE /opt/image-previewer/run-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
# COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
    -ldflags "$LDFLAGS" \
    -o ${BIN_FILE} cmd/image-previewer/*

FROM alpine:3.9

LABEL SERVICE="image-previewer"
LABEL MAINTAINERS="gromov.91@gmail.com"

ENV BIN_FILE "/opt/image-previewer/run-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
