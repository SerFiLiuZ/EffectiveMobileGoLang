CREATE TABLE IF NOT EXISTS people (
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100),
    PRIMARY KEY (name, surname)
);
