FROM golang:1.23 AS build
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -ldflags "-s" -o fuwa-server .

FROM scratch
WORKDIR /
COPY --from=build /app/fuwa-server .
EXPOSE 6942
CMD ["./fuwa-server"]
