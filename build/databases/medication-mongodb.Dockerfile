FROM mongo:8.2.1

COPY build/databases/mongo-init.js /docker-entrypoint-initdb.d/

EXPOSE 27017