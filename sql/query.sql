CREATE TABLE files
(
   	Id serial PRIMARY KEY,
    Name varchar(255),
    FileId varchar(255),
    FileUniqueId varchar(255),
    FileSize integer DEFAULT 0
);

CREATE TABLE directories
(
   	Id serial PRIMARY KEY,
    Name varchar(255),
    UserID integer,
    Files integer[],
    Directories integer[],
    Size integer DEFAULT 0
);

CREATE TABLE users
(
   	Id serial PRIMARY KEY,
    Username varchar(255),
    UserID integer,
    FirstName varchar(255),
    LastName varchar(255),
    CurrentDirectory integer REFERENCES directories (Id)
);