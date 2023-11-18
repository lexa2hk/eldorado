FROM postgres:alpine
LABEL Name="Eldorado Database"
LABEL Version="0.1"
COPY ./db/migrations/*.up.sql /docker-entrypoint-initdb.d/
COPY ./db/migrations/init.sql /docker-entrypoint-initdb.d/
