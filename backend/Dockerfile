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

FROM scratch

COPY --from=build /out/calculator-api /calculator-api

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/calculator-api"]
