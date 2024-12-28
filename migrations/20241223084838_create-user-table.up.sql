CREATE TABLE users (
    id INT GENERATED ALWAYS AS IDENTITY,
    username TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    
    PRIMARY KEY (id)
);