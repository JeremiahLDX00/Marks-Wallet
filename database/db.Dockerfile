FROM mysql:8.0.23

EXPOSE 9073

COPY ./*.sql /docker-entrypoint-initdb.d/