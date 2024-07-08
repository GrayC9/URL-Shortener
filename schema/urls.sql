CREATE DATABASE IF NOT EXISTS urlsh;
USE urlsh;
CREATE TABLE urls (
                      id INT AUTO_INCREMENT PRIMARY KEY,
                      short_code VARCHAR(255) NOT NULL,
                      original_url TEXT NOT NULL,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);