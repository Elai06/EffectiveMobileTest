CREATE TABLE IF NOT EXISTS people (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    surname VARCHAR(255),
    age INT,
    gender VARCHAR(10),
    nationality VARCHAR(50)
);