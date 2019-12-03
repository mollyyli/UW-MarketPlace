create table if not exists Users (
    ID int not null auto_increment primary key,
    Email varchar(128) not null,
    PassHash varchar(64) not null,
    Username varchar(255) not null,
    FirstName varchar(255) not null,
    LastName varchar(255) not null,
    PhotoURL varchar(300) not null
);

create table if not exists UserSignIns (
    ID int not null auto_increment primary key,
    UserID int not null,
    SignInTime varchar(255) not null,
    IP varchar(255) not null
);