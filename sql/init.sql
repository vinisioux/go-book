-- Active: 1710614295246@@127.0.0.1@5432@gobook
DO $$ 
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'gobook') THEN
    CREATE DATABASE gobook;
  END IF;
END $$;

\c gobook;

DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS users;

CREATE TABLE users(
  id serial primary key,
  name varchar(40) not null,
  nickname varchar(20) not null unique,
  email varchar(50) not null unique,
  password varchar(150) not null,
  created_at timestamp default current_timestamp
);

CREATE TABLE followers(
  user_id int not null,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  follower_id int not null,
  FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, follower_id)
);

CREATE TABLE posts(
  id serial PRIMARY KEY,
  title VARCHAR(50) not null,
  content VARCHAR(150) not null,
  author_id int not null,
  Foreign Key (author_id) REFERENCES users(id) ON DELETE CASCADE,
  likes int DEFAULT 0,
  created_at timestamp default current_timestamp
);
