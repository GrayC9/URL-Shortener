CREATE DATABASE IF NOT EXISTS urlsh;
USE urlsh;
ALTER TABLE urls ADD COLUMN click_count INT DEFAULT 0;
ALTER TABLE urls ADD COLUMN last_accessed_at TIMESTAMP NULL;
CREATE TABLE urls (
                      id INT AUTO_INCREMENT PRIMARY KEY,
                      short_code VARCHAR(255) NOT NULL,
                      original_url TEXT NOT NULL,
                      click_count INT DEFAULT 0,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
