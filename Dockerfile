FROM golang:alpine as build

ENV PROJECT_PATH="/go/src/url-shortener-user-service"

WORKDIR $PROJECT_PATH

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
RUN go mod download

COPY . $PROJECT_PATH

RUN go build -ldflags="-s -w" -o /app ./cmd


FROM golang:alpine

RUN apk --no-cache add curl

ENV WAITFORIT_VERSION="v2.4.1"
RUN curl -o /usr/local/bin/waitforit -sSL https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 && \
    chmod +x /usr/local/bin/waitforit

ARG PORT
ARG LOG_LEVEL
ARG MONGO_URI
ARG MONGO_DB
ARG MONGO_COLLECTION
ARG REDIS_URL

ENV PORT=$PORT
ENV LOG_LEVEL=$LOG_LEVEL
ENV MONGO_URI=$MONGO_URI
ENV MONGO_DB=$MONGO_DB
ENV MONGO_COLLECTION=$MONGO_COLLECTION
ENV REDIS_URL=$REDIS_URL

COPY --from=build /app /app

EXPOSE ${PORT}

CMD ["/app"]
