FROM node:lts-alpine

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm ci --only=production

COPY dist/database/ ./

CMD [ "node", "migration.js" ]
