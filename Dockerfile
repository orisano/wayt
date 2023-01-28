FROM golang:1.19 AS build
WORKDIR /app
COPY . .
RUN go build .
ENTRYPOINT /app/wayt

FROM scratch
COPY --from=build /app/wayt /bin/wayt
ENTRYPOINT /bin/wayt
