FROM node:8

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY ./package.json /usr/src/app
COPY ./package-lock.json /usr/src/app
RUN npm install

COPY . /usr/src/app/
COPY ./env.prod.sh /usr/src/app/env.prod.sh
RUN npm run build

RUN chmod -R 777 /usr/src/app/dist/

FROM abiosoft/caddy:no-stats

COPY --from=0 /usr/src/app/dist/ /srv
COPY Caddyfile /etc/