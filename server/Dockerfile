ARG GO_VERSION=1.23.0
FROM golang:${GO_VERSION}-alpine AS build

COPY go.mod go.sum ./
RUN go mod download -x

COPY server/ server/
COPY core/ core/
RUN go build -o /minioserve ./server

FROM alpine:3.12

COPY --from=build /minioserve /usr/local/bin/minioserve

ENV MINICHAT_PORT=3000
EXPOSE $MINICHAT_PORT

CMD [ "minioserve" ]