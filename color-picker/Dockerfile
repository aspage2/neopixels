FROM node:22.12

WORKDIR /app

COPY yarn.lock package.json package-lock.json ./

RUN yarn

COPY . ./

EXPOSE 5713

ENTRYPOINT ["yarn", "dev"]

