DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users(
    "id" VARCHAR(32) NOT NULL UNIQUE,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL,

    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "created_by" VARCHAR(32) NOT NULL,
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_by" VARCHAR(32) NOT NULL,
    "deleted_at" TIMESTAMP,
    "deleted_by" VARCHAR(32),

    PRIMARY KEY (id)
);
DROP TABLE IF EXISTS users_properties_changes_history;
CREATE TABLE IF NOT EXISTS users_properties_changes_history(
    "user_id" VARCHAR(32) NOT NULL,
    "name" VARCHAR(32) NOT NULL,
    "changed_from" VARCHAR(255),
    "changed_to" VARCHAR(255) NOT NULL,

    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "created_by" VARCHAR(32) NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id)
);