CREATE DATABASE IF NOT EXISTS gobook;
USE gobook;

DROP TABLE IF EXISTS users;

CREATE TABLE users(
  id serial primary key,
  name varchar(40) not null,
  nickname varchar(20) not null unique,
  email varchar(50) not null unique,
  password varchar(40) not null,
  created_at timestamp default current_timestamp
);