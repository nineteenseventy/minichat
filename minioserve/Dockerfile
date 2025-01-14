ARG GO_VERSION=1.23.0
FROM golang:${GO_VERSION}-alpine AS build

COPY go.mod go.sum ./
RUN go mod download -x

COPY core/ core/
COPY minioserve/ minioserve/
RUN go build -o /minioserve ./minioserve

FROM alpine:3.12

COPY --from=build /minioserve /usr/local/bin/minioserve

ENV MINIOSERVE_HOST=0.0.0.0
ENV MINIOSERVE_PORT=3000
EXPOSE $MINIOSERVE_PORT

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "wget", "-q", "http://localhost:3000/health" ]

CMD [ "minioserve" ]