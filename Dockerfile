ARG APP_NAME=consularis
ARG APP_PATH=/$APP_NAME
ARG SRC_PATH=/go/src/github.com/leosunmo/$APP_NAME

# Build container
FROM golang:1.10.2 AS builder

# Import args
ARG APP_NAME
ARG SRC_PATH

COPY . $SRC_PATH
WORKDIR $SRC_PATH

RUN CGO_ENABLED=0 go build

# Build final image with binary
FROM alpine:3.7

# Import args
ARG APP_NAME
ARG SRC_PATH

COPY --from=builder $SRC_PATH/$APP_NAME /$APP_NAME

CMD ["/consularis"]
