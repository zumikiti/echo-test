FROM golang:1.18.3-alpine as build

WORKDIR /go/app

COPY . .

RUN apk add --no-cache git \
 && go build -o app

FROM alpine

WORKDIR /app

COPY --from=build /go/app/app .

RUN addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /app/app

CMD ["./app"]
