FROM node:11 AS build-frontend

WORKDIR /tmp

ADD frontend/*.json ./
ADD frontend/*.lock ./
RUN yarn install

ADD frontend/. .
RUN yarn run lint
RUN yarn run build --prod

FROM golang AS build-server

WORKDIR /build

RUN apt update && apt -y install bzip2

# copy SDE version and download EVE SDE
COPY sde.* ./
RUN ./sde.sh

# copy dependency information and fetch them
COPY go.mod ./
RUN go mod download

# copy sources
COPY . .

# build and install (without C-support, otherwise there issue because of the musl glibc replacement on Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -a cmd/server/server.go

FROM alpine
# update CA certificates
RUN apk update && apk add ca-certificates
WORKDIR /usr/titan
COPY --from=build-frontend /tmp/dist ./frontend/dist
COPY --from=build-server /build/server .
COPY --from=build-server /build/sde.version .
COPY --from=build-server /build/sde-* .
COPY restore.sh .
CMD ["./docker-entrypoint.sh"]
