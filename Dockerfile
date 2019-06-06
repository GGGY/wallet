FROM golang:1.12-alpine as builder

WORKDIR /project

RUN set -xe  && \
    apk update && apk upgrade  && \
    apk add --no-cache make git curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v3.5.2/migrate.linux-amd64.tar.gz | tar xvz && \
    cp migrate.linux-amd64 /migrate

COPY . .

RUN make dep && \
      make build && \
      cp build/wallet /service


FROM scratch

COPY --from=builder /service /service
COPY --from=builder /migrate /migrate


CMD ["/service"]