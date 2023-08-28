CREATE TABLE users 
(
    userID VARCHAR(255) NOT NULL,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE notes 
(
    userID VARCHAR(255) NOT NULL,
    username VARCHAR(50) NOT NULL,
    note TEXT NOT NULL
);