-- +goose Up
-- +goose StatementBegin
-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO users (name) VALUES
('Barmaley'),
('Tyanitolkai'),
('Hippopotamus'),
('Monkey Chichi'),
('Dog Avva'),
('Parrot Karudo'),
('Kangaroo'),
('Lion'),
('Crocodile'),
('Giraffe');

-- Create scheduled table
CREATE TABLE scheduled (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    cure_name VARCHAR(255) NOT NULL,
    frequency BIGINT NOT NULL,
    duration BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scheduled;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
