-- +goose Up
-- +goose StatementBegin
-- Add extenstions
CREATE EXTENSION IF NOT EXISTS "pgcrypto"; 
CREATE EXTENSION IF NOT EXISTS postgis;

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
    duration_days INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO scheduled (user_id, cure_name, doses_per_day, duration_days, created_at) VALUES
(1, 'Paracetamol', 10, 10, NOW()),  
(2, 'Ibuprofen', 2, 5, NOW()),    
(3, 'Aspirin', 15, 14, NOW()),     
(4, 'Amoxicillin', 3, 12, NOW()),  
(5, 'Ozempic', 2, 240, NOW()),   
(6, 'Vitamin D', 1, 0, NOW()),                  
(7, 'Red pill', 1, 0, NOW()),                 
(8, 'Atorvastatin', 1, 7, NOW()), 
(9, 'Cetirizine', 2, 10, NOW()),    
(10, 'Omeprazole', 1, 10, NOW());  

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scheduled;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
