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
    doses_per_day INT NOT NULL,
    duration BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO scheduled (user_id, cure_name, doses_per_day, duration, created_at) VALUES
(1, 'Paracetamol', 3, 1209600000000000, NOW()),  -- 2 weeks
(2, 'Ibuprofen', 2, 2419200000000000, NOW()),    -- 4 weeks
(3, 'Aspirin', 4, 31536000000000000, NOW()),     -- 1 year
(4, 'Amoxicillin', 3, 1209600000000000, NOW()),  -- 2 weeks
(5, 'Ozempic', 2, 31536000000000000, NOW()),   -- 1 year
(6, 'Vitamin D', 1, 0, NOW()),                   -- Permanent use
(7, 'Lisinopril', 1, 0, NOW()),                  -- Permanent use
(8, 'Atorvastatin', 1, 31536000000000000, NOW()), -- 1 year
(9, 'Cetirizine', 2, 604800000000000, NOW()),    -- 1 week
(10, 'Omeprazole', 1, 1209600000000000, NOW());  -- 2 weeks

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scheduled;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
