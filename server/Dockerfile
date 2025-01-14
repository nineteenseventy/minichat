ARG GO_VERSION=1.23.0
FROM golang:${GO_VERSION}-alpine AS build

COPY go.mod go.sum ./
RUN go mod download -x

COPY server/ server/
COPY core/ core/
RUN go build -o /server ./server

FROM alpine:3.12

COPY --from=build /server /usr/local/bin/server

ENV MINICHAT_HOST=0.0.0.0
ENV MINICHAT_PORT=3000
EXPOSE $MINICHAT_PORT

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "wget", "-q", "http://localhost:3000/health" ]

CMD [ "server" ]