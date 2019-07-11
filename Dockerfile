FROM golang:latest

# RUN apk update && apk add gcc libc-dev make git

RUN mkdir /app
WORKDIR /app
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo .

WORKDIR /app/plugins
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux make
RUN cp -r build ../
WORKDIR /app
# CMD ["./supergroup.mixin.one"]
EXPOSE 7001
EXPOSE 9001

# FROM alpine:3.4
# COPY --from=0 /app/supergroup.mixin.one /app/
# RUN mkdir /app/plugins
# COPY --from=0 /app/plugins/build /app/plugins
# WORKDIR /app
# COPY ./config.yaml /app/config.yaml
# RUN apk update \
#         && apk upgrade \
#         && apk add --no-cache \
#         ca-certificates \
#         && update-ca-certificates 2>/dev/null || true
# EXPOSE 7001
# EXPOSE 9001

CMD ["./supergroup.mixin.one"]