DO $$ 
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'gobook') THEN
    CREATE DATABASE gobook;
  END IF;
END $$;

\c gobook;

DROP TABLE IF EXISTS users;

CREATE TABLE users(
  id serial primary key,
  name varchar(40) not null,
  nickname varchar(20) not null unique,
  email varchar(50) not null unique,
  password varchar(40) not null,
  created_at timestamp default current_timestamp
);