CREATE TABLE IF NOT EXISTS files
(
   	Id serial PRIMARY KEY,
    UserID integer,
    Name varchar(255),
    FileId varchar(255),
    FileUniqueId varchar(255),
    FileSize integer DEFAULT 0,
    FileType varchar(255),
    Created timestamp,
    ThumbnailFileId varchar(255) DEFAULT '',
    ThumbnailSource varchar(255) DEFAULT '',
    FileSource varchar(255) DEFAULT '',
    SharedId varchar(255) NOT NULL DEFAULT gen_random_uuid(),
    IsShared bool DEFAULT false
);

CREATE TABLE IF NOT EXISTS directories
(
   	Id serial PRIMARY KEY,
    ParentId integer DEFAULT -1,
    Name varchar(255),
    UserID integer,
    Files integer[] DEFAULT '{}',
    Directories integer[] DEFAULT '{}',
    Size integer DEFAULT 0,
    Path varchar(255),
    Created timestamp
);

CREATE TABLE IF NOT EXISTS users
(
   	Id serial PRIMARY KEY,
    Username varchar(255),
    UserID integer,
    FirstName varchar(255),
    LastName varchar(255),
    CurrentDirectory integer REFERENCES directories (Id),
    Hash varchar(255)
);
