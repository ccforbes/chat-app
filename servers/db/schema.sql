create table if not exists Users (
    id int not null auto_increment primary key,
    email varchar(345) not null,
    passHash varbinary(128) not null,
    username varchar(255) not null,
    firstName varchar(64) not null,
    lastName varchar(64) not null,
    photoUrl varchar(128) not null
    unique key unique_email (email),
    unique key unique_username (username)
);