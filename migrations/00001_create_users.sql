-- +goose Up
CREATE TABLE users (
	id VARCHAR(255) PRIMARY KEY,
	username VARCHAR(50) UNIQUE NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	hashed_password VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;