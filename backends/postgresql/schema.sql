-- Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

CREATE TABLE users (
    id bigserial primary key,
    isadmin bool NOT NULL default false,
    isbanned bool NOT NULL default false,
    username varchar(30) UNIQUE NOT NULL,
    email varchar(30) UNIQUE NOT NULL,
    bio varchar(200) NOT NULL default '',
    lastlogin timestamp default current_timestamp,
    createdon timestamp default current_timestamp,
    lastupdated timestamp default current_timestamp
);

CREATE TABLE channels (
    id bigserial primary key,
    ownerid bigserial references users(id),
    dmid bigserial references users(id),
    name varchar(30) UNIQUE NOT NULL,
    description text NOT NULL default '',
    topic text NOT NULL default '',
    isprivate boolean NOT NULL default false,
    createdon timestamp default current_timestamp,
    lastupdated timestamp default current_timestamp
);

CREATE TABLE characters (
    id bigserial primary key,
    userid bigserial references users(id),
    channelid bigserial references channels(id),
    name varchar(30) NOT NULL default '',
    description text NOT NULL default '',
    createdon timestamp default current_timestamp,
    lastupdated timestamp default current_timestamp,
    UNIQUE (name, channelid),
    UNIQUE (userid, channelid)
);

CREATE TABLE messages (
    id bigserial primary key,
    characterid bigserial references characters(id),
    channelid bigserial references channels(id),
    content varchar(200) NOT NULL,
    isstory boolean NOT NULL,
    createdon timestamp default current_timestamp,
    lastupdated timestamp default current_timestamp
);

-- Use function provided from functions.sql to handle updating lastmodified timestamps on updates
CREATE TRIGGER users_updated_at_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();
CREATE TRIGGER channels_updated_at_modtime BEFORE UPDATE ON channels FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();
CREATE TRIGGER characters_updated_at_modtime BEFORE UPDATE ON characters FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();
CREATE TRIGGER messages_updated_at_modtime BEFORE UPDATE ON messages FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();

-- Sample data
INSERT INTO users
(username, email, isadmin, isbanned) VALUES
('andrew.w.boutin@gmail.com', 'andrew.w.boutin@gmail.com', true, false),
('banneduser', 'banneduser@fake.com', false, true),
('adminuser', 'adminuser@fake.com', true, false),
('regularuser', 'regularuser@fake.com', false, false);

INSERT INTO channels
(ownerid, dmid, name, description, topic, isprivate) VALUES
(1, 1, 'my public channel', 'my public channel description', 'some topic', false),
(1, 1, 'my private channel', 'my private channel description', '', true);

INSERT INTO characters
(userid, channelid, name, description) VALUES
(1, 1, 'character1', 'character1...'),
(1, 2, 'character2', 'character2...');

INSERT INTO messages
(characterid, channelid, content, isStory) VALUES
(1, 1, 'message one story public channel', true),
(1, 1, 'messsage two meta public channel', false),
(1, 2, 'message one story private channel', true),
(1, 2, 'messsage two meta private channel', false);
