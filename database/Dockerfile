FROM node:lts-alpine

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm ci --include=dev

COPY database/ ./

CMD [ "npx", "tsx", "migration.ts" ]
