FROM node:lts-alpine

COPY .output /app

WORKDIR /app

ENV NITRO_PORT=3000
ENV NITRO_HOST=0.0.0.0

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "curl", "-f", "http://localhost:${NITRO_PORT}/api/health" ]

EXPOSE $NITRO_PORT

CMD [ "node", "./server/index.mjs" ]