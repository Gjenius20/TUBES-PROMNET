FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server cmd/server/main.go

FROM gcr.io/distroless/static
COPY --from=build /server /server
ENV GIN_MODE=release
EXPOSE 8080
CMD ["/server"]
