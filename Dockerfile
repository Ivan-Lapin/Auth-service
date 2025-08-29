FROM golang:1.24.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth_service ./service/cmd/auth/main.go

FROM golang:1.24.5

WORKDIR /root/

COPY --from=builder /app/auth_service .

COPY /service/migrations /root/migrations

EXPOSE 8877

CMD [ "./auth_service" ]