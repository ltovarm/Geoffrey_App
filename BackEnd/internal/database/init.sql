CREATE DATABASE temp;

\c temp;

CREATE TABLE temperatures (
    id SERIAL PRIMARY KEY,
    temperature FLOAT NOT NULL,
    date TIMESTAMP NOT NULL
);
