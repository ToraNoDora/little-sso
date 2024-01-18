FROM golang:1.19-alpine as builder

ENV PROXY=https://proxy-test/repository/go-google-proxy/,direct \
    PLATFORM=linux \
    OUTPUT_PATH=tmp/bin \
    BINARY_NAME=little_sso \
    MAIN_PACKAGE_PATH=./cmd/sso/main.go \
    BUILD_PATH=/go/src/github.com/ToraNoDora/little-sso \
    ARCH=amd64

WORKDIR ${BUILD_PATH}

COPY ./sso .

RUN apk update && apk add --no-cache curl gcc git libc-dev

# revive (go lint successor)
RUN go install github.com/mgechev/revive@latest && \
    revive ./...

# gosec - Golang Security Checker
# RUN go install github.com/securego/gosec/v2/cmd/gosec@latest && \
#     gosec ./...

RUN go env -w GOSUMDB=off \
    && go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOOS=${PLATFORM} \
    # go env -w GOPROXY=${PROXY} \
    && GOARCH=${ARCH}

# Get go dependencies
RUN go mod download

RUN go build -o=${OUTPUT_PATH}/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

RUN cp ${BUILD_PATH}/${OUTPUT_PATH}/${BINARY_NAME} /bin/${BINARY_NAME}
RUN cp ${BUILD_PATH}/configs/config.locale.yaml /bin/config.locale.yaml


FROM alpine:3

RUN apk update && apk add --no-cache ca-certificates tzdata libc6-compat

COPY --from=builder /bin/little_sso /little_sso
COPY --from=builder /bin/config.locale.yaml configs/config.locale.yaml

EXPOSE 44044

CMD ["/little_sso"]
