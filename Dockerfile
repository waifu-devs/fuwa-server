FROM golang:1.23 AS build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . ./

RUN go build -v -ldflags "-s -X 'main.Version=${GIT_BRANCH}.${GIT_HASH}'" -o fuwa-server

FROM scratch

COPY --from=build /app/fuwa-server /fuwa-server

EXPOSE 6942
ENTRYPOINT ["/fuwa-server"]
