CREATE DATABASE house;

\c house;

CREATE TABLE temperatures (
    id SERIAL PRIMARY KEY,
    data JSONB
);
