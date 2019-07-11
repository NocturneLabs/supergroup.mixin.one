FROM node:8

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY ./package.json /usr/src/app
RUN npm install

COPY . /usr/src/app/
COPY ./env.prod.sh /usr/src/app/env.prod.sh
RUN npm run build

RUN chmod -R 777 /usr/src/app/dist/

VOLUME /usr/src/app/dist

FROM nginx:1.16.0

COPY --from=0 /usr/src/app/dist/ /var/www/mixin_group/dist
COPY ./web.conf  /etc/nginx/conf.d/
COPY ./api.conf  /etc/nginx/conf.d/