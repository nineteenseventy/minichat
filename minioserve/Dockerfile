ARG GO_VERSION=1.23.0
FROM golang:${GO_VERSION}-alpine AS build

COPY go.mod go.sum ./
RUN go mod download -x

COPY main.go .
RUN go build -o /minioserve

FROM alpine:3.12

COPY --from=build /minioserve /usr/local/bin/minioserve

ENV MINIOSERVE_PORT=3000
EXPOSE $MINIOSERVE_PORT

CMD [ "minioserve" ]