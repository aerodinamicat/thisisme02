DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS persons;
CREATE TABLE IF NOT EXISTS persons(
    first_name VARCHAR(255) NOT NULL,
    second_name VARCHAR(255),
    first_surname VARCHAR(255) NOT NULL,
    second_surname VARCHAR(255),
    gender VARCHAR(32) NOT NULL,
    birth_date TIMESTAMP NOT NULL,

    user_id VARCHAR(32) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY(user_id) REFERENCES users(id)
);