FROM golang:1.16.12-alpine3.15 AS go-builder

WORKDIR /usr/src/app

RUN apk add --update \
        curl \
        gcc \
        git \
        make \
        musl-dev \
        tzdata

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /app


FROM alpine:3.15

# Install packages required by the image
RUN apk add --update \
        bash \
        ca-certificates \
        coreutils \
        curl \
        jq \
        tzdata \
        openssl \
    && rm /var/cache/apk/*

COPY --from=go-builder /app ./

COPY . .

CMD [ "./app" ]