FROM alpine

# update CA certificates
RUN apk update && apk add ca-certificates postgresql-client
WORKDIR /usr/titan
RUN mkdir -p frontend/dist/titan-frontend
COPY frontend/dist/titan-frontend ./frontend/dist/titan-frontend/
COPY server  .
COPY sde.version .
COPY sde-* ./
ADD restore.sh .
ADD docker-entrypoint.sh .
ADD sql sql
CMD ["./docker-entrypoint.sh"]
