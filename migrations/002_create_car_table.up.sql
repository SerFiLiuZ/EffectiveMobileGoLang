CREATE TABLE IF NOT EXISTS car (
    regNum VARCHAR(20) PRIMARY KEY,
    mark VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER,
    owner_name VARCHAR(100) NOT NULL,
    owner_surname VARCHAR(100) NOT NULL,
    owner_patronymic VARCHAR(100),
    CONSTRAINT fk_owner FOREIGN KEY (owner_name, owner_surname, owner_patronymic) REFERENCES people(name, surname, patronymic)
);