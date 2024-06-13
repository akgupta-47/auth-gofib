# FROM golang:1.18.10-alpine3.17

# WORKDIR /app

# COPY . .

# RUN go build -o main main.go

# EXPOSE 5000

# CMD [ "/app/main" ]
# ----------------------------------------
# Multi Stage build file to reduce size by only having the go binary in the image

# Build Stage
FROM golang:1.18.10-alpine3.17 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

# Run Stage
FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .

EXPOSE 5000

CMD [ "/app/main" ]
