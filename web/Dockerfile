FROM node:16.13.1 AS build

WORKDIR /usr/app

RUN curl -f https://get.pnpm.io/v6.16.js | node - add --global pnpm

COPY pnpm-lock.yaml ./

RUN pnpm fetch

COPY . .

RUN pnpm install -r --offline
RUN pnpm run build

FROM nginx:alpine
COPY --from=build /usr/app/dist/ /usr/share/nginx/html/
COPY ./docker/nginx/nginx.conf /etc/nginx/conf.d/
