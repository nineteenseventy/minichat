FROM node:lts-alpine

COPY .output /app

WORKDIR /app

CMD [ "node", "./server/index.mjs" ]