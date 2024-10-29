FROM golang:1.23-alpine AS build
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN apk add --no-cache \
	# Important: required for go-sqlite3
	gcc \
	# Required for Alpine
	musl-dev
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build -v -ldflags "-s -extldflags '-static'" -o fuwa-server .

FROM scratch
WORKDIR /
COPY --from=build /app/fuwa-server .
EXPOSE 6942
CMD ["./fuwa-server"]
