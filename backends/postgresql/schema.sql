CREATE TABLE users (
    id bigserial primary key,
    name varchar(30) NOT NULL,
    bio varchar(200) NOT NULL,
    createdon timestamp default current_timestamp,
    lastupdated timestamp default current_timestamp
);

CREATE TABLE channels (
    id bigserial primary key,
    ownerid bigserial references users(id),
    dmid bigserial references users(id),
    name varchar(30) NOT NULL,
    description text NOT NULL,
    isprivate boolean NOT NULL,
    createdon timestamp default current_timestamp,
    lastupdated timestamp default current_timestamp
);

CREATE TABLE messages (
    id bigserial primary key,
    userid bigserial references users(id),
    content varchar(200) NOT NULL,
    createdon timestamp default current_timestamp
);

INSERT INTO users
(name, bio) VALUES
('user1', 'my first user');

INSERT INTO channels
(ownerid, dmid, name, description, isprivate) VALUES
(1, 1, 'my channel', 'my channel description', false);