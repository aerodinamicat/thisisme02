FROM postgres:10.3

COPY initial.sql /docker-entrypoint-initdb.d/1.sql

EXPOSE 5432

CMD ["postgres"]