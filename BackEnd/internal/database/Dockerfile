# Official image of PostgreSQL
FROM postgres:latest

# Define the enviroment variable for PostgreSQL settings
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD mysecretpassword
ENV POSTGRES_DB mydb

# Copy initialisation SQL script in the image
COPY init.sql /docker-entrypoint-initdb.d/

# Expose the default port of PostgreSQL
EXPOSE 5432
