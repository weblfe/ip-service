FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build/zero

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/ipservice ./api/ipservice.go


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
#ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/ipservice /app/ipservice

EXPOSE 8080
VOLUME /app/tmp/db

CMD ["./ipservice"]
