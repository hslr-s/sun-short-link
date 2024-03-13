
# build backend
FROM golang:1.21-alpine3.18 as server_image

WORKDIR /build

COPY . .

# 中国国内源
RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories \
    && go env -w GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache bash curl gcc git musl-dev

RUN go env -w GO111MODULE=on \
    && go build -o sun-short-link main.go



# run_image
FROM alpine

WORKDIR /app

COPY --from=server_image /build/sun-short-link /app/sun-short-link

# 中国国内源
RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories

EXPOSE 3002

RUN apk add --no-cache bash ca-certificates su-exec tzdata \
    && chmod +x ./sun-short-link \
    && ./sun-short-link -i -c ./sun-short-link.yml

CMD ./sun-short-link -c sun-short-link.yml
