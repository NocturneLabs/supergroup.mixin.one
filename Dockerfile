FROM golang:alpine
RUN apk update && apk add --no-cache git gcc g++ make

RUN mkdir /app
WORKDIR /app
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN mkdir /app/plugins
COPY plugins/go.mod /app/plugins
COPY plugins/go.sum /app/plugins
RUN cd /app/plugins && go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .
RUN cd /app/plugins && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 make

FROM alpine:3.4
COPY --from=0 /app/supergroup.mixin.one /app/
RUN mkdir /app/plugins
COPY --from=0 /app/plugins/build /app/plugins
WORKDIR /app
COPY ./config.yaml /app/config.yaml
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
EXPOSE 7001
EXPOSE 9001

CMD ["./supergroup.mixin.one"]