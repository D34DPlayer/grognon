FROM golang:1.24-alpine AS builder-go
RUN apk add --no-cache --update gcc g++

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/grognon ./main.go

FROM node:22-alpine AS deps-node
WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM deps-node AS builder-node
COPY . .

RUN pnpm build

FROM deps-node AS runner-node
COPY --from=builder-node /app/bootstrap /app/bootstrap

CMD ["/app/bootstrap/assets/ssr.js"]

FROM alpine:latest AS runner-go
WORKDIR /app
RUN cd /app

COPY ./index.html /app/index.html
COPY --from=builder-go /app/bin/grognon /app/bin/grognon
COPY --from=builder-node /app/public /app/public

ENTRYPOINT ["/app/bin/grognon"]
