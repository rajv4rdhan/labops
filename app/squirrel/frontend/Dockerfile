FROM golang:1.25 AS build

WORKDIR /src

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

FROM scratch

WORKDIR /app

COPY --from=build /out/server /app/server
COPY index.html /app/index.html
COPY styles.css /app/styles.css
COPY app.js /app/app.js

EXPOSE 8080

ENTRYPOINT ["/app/server"]
