<h1 align="center">Hi, this is a service, url_shortener <a a </a> 
<img src="https://github.com/blackcater/blackcater/raw/main/images/Hi.gif" height="32"/></h1>
<h3 align="center"> Welcome to the URL shortening service documentation. I created a REST API that shortens URLs.</h3>

<h3 align="left" style="font-size: 1.5em;"> Description of how my URL shortening service works with authentication </h3>
<p align="left">URL shortening:

The first time a user shortens a link, a new entry is created in the database that includes the original URL (original_url) and the generated short URL (short_url).
Information about the user who created the shortcut is also stored in the database if the user is authenticated.
At the same time, a counter is updated showing how many times this original_url has been shortened.

Redirect to short URL:

When someone navigates to a short URL (short_url), the service looks up the corresponding original_url in the database.
After finding the original URL, the service redirects the user to this URL.
The database is updated with a hop count showing how many times this short_url has been used.

Caching:

To speed up query processing, popular URLs and their shortenings are cached in memory.
With each request, the cache is first checked, and only if there is no data in the cache, the database is accessed.</p>

<h3 align="left" style="font-size: 1.5em;"> Description of how the cache works </h3>
<p align="left">Cache :
  
The URL shortening service implements caching to improve performance and speed up request processing. Caching is used to store popular URLs and their abbreviations in RAM, thereby avoiding frequent database calls. Key aspects of caching include warming up the cache, limiting its size, and using the LRU (Least Recently Used) strategy.


Here is the work of the cache and its interaction with the database:
![image](https://github.com/user-attachments/assets/d1ef5ab0-b962-4cdf-b64c-000fcc29cd81)

Using the LRU Strategy:
The LRU (Least Recently Used) caching strategy is used to manage cache contents. This strategy ensures that the least recently used items are removed when new data is added, allowing for efficient use of limited cache space.

Maintain usage order: The cache keeps track of the last time each item was used. This allows you to track which items have been used recently and which have not.
Removing old items: When it is necessary to make room for a new item, the cache removes the least recently used item, that is, the one that was used the longest ago.

How LRU cache works:

![image](https://github.com/user-attachments/assets/edc65c19-46e2-42f2-93b5-09e702269073)


Deploying the application to the server:

To deploy the application, you need to run my Docker container on your server.

docker run -d -p 8943:8080 --env DB_DSN="user:password@tcp(db_host:db_port)/db_name" grayc9/shortener:1722964983

You need to pick up your database and create the following tables in it:

# Table
```shell
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
);  "
```

</p>

# Configuration 
```shell
export DB_DSN="user:password@tcp(db_host:db_port)/db_name"
export JWT_SECRET="your_jwtkey"
```
