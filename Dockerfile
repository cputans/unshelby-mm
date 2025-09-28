FROM golang:1.24-alpine AS build

ENV GOOS linux

WORKDIR /app
RUN apk update && apk add --update gcc musl-dev alpine-sdk
ADD . .
RUN go build -o unshelbymm *.go 

FROM alpine
COPY --from=build /app/unshelbymm .
CMD ["/app/unshelbymm"]