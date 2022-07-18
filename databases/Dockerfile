FROM postgres:latest

COPY initial.sql /docker-entrypoint-initdb.d/1.sql

EXPOSE 5432

CMD ["postgres"]