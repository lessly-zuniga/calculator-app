FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/calculator-api \
    ./cmd/server

FROM alpine:3.23

RUN addgroup --system app && adduser --system --ingroup app app

COPY --from=build /out/calculator-api /calculator-api

USER app

EXPOSE 8080

ENTRYPOINT ["/calculator-api"]
