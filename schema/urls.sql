CREATE DATABASE IF NOT EXISTS urlsh;
USE urlsh;
ALTER TABLE urls ADD COLUMN click_count INT DEFAULT 0;
ALTER TABLE urls ADD COLUMN last_accessed_at TIMESTAMP NULL;
ALTER TABLE urls DROP COLUMN last_accessed_at;
CREATE TABLE urls (
                      id INT AUTO_INCREMENT PRIMARY KEY,
                      short_code VARCHAR(255) NOT NULL,
                      original_url TEXT NOT NULL,
                      click_count INT DEFAULT 0,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
                      id INT AUTO_INCREMENT PRIMARY KEY,
                      username VARCHAR(255) NOT NULL UNIQUE,
                      hash_password VARCHAR(255) NOT NULL,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);  
DROP TABLE users;

SELECT * FROM users;
