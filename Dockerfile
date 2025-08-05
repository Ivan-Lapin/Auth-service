FROM golang:latest AS builder

WORKDIR /app

COPY go.mode go sum ./

RUN go mod download

COPY . .

RUN go biuld -o auth-service /service/cmd/auth

FROM golang:latest

COPY --from=builder /app/service .

COPY ./service/internal/config/config.yaml ./service/internal/config/config.yaml

EXPOSE 8878

ENTRYPOINT [ "./auth-service" ]

