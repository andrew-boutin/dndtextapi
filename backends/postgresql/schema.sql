-- Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

CREATE TABLE users (
    id bigserial primary key,
    is_admin bool NOT NULL default false,
    is_banned bool NOT NULL default false,
    username varchar(30) UNIQUE NOT NULL,
    email varchar(30) UNIQUE NOT NULL,
    bio varchar(200) NOT NULL default '',
    last_login timestamp default current_timestamp,
    created_on timestamp default current_timestamp,
    last_updated timestamp default current_timestamp
);

CREATE TABLE channels (
    id bigserial primary key,
    owner_id bigserial references users(id),
    dm_id bigserial references users(id),
    name varchar(30) UNIQUE NOT NULL,
    description text NOT NULL default '',
    topic text NOT NULL default '',
    is_private boolean NOT NULL default false,
    created_on timestamp default current_timestamp,
    last_updated timestamp default current_timestamp
);

CREATE TABLE characters (
    id bigserial primary key,
    user_id bigserial references users(id),
    channel_id bigserial references channels(id),
    name varchar(30) NOT NULL default '',
    description text NOT NULL default '',
    created_on timestamp default current_timestamp,
    last_updated timestamp default current_timestamp,
    UNIQUE (name, channel_id),
    UNIQUE (user_id, channel_id)
);

CREATE TABLE messages (
    id bigserial primary key,
    character_id bigserial references characters(id),
    channel_id bigserial references channels(id),
    content varchar(200) NOT NULL,
    is_story boolean NOT NULL,
    created_on timestamp default current_timestamp,
    last_updated timestamp default current_timestamp
);

-- Use function provided from functions.sql to handle updating lastmodified timestamps on updates
CREATE TRIGGER users_updated_at_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();
CREATE TRIGGER channels_updated_at_modtime BEFORE UPDATE ON channels FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();
CREATE TRIGGER characters_updated_at_modtime BEFORE UPDATE ON characters FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();
CREATE TRIGGER messages_updated_at_modtime BEFORE UPDATE ON messages FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();

-- Sample data
INSERT INTO users
(username, email, is_admin, is_banned) VALUES
('andrew.w.boutin@gmail.com', 'andrew.w.boutin@gmail.com', true, false),
('banneduser', 'banneduser@fake.com', false, true),
('adminuser', 'adminuser@fake.com', true, false),
('regularuser', 'regularuser@fake.com', false, false);

INSERT INTO channels
(owner_id, dm_id, name, description, topic, is_private) VALUES
(1, 1, 'my public channel', 'my public channel description', 'some topic', false),
(1, 1, 'my private channel', 'my private channel description', '', true);

INSERT INTO characters
(user_id, channel_id, name, description) VALUES
(1, 1, 'character1', 'character1...'),
(1, 2, 'character2', 'character2...');

INSERT INTO messages
(character_id, channel_id, content, is_story) VALUES
(1, 1, 'message one story public channel', true),
(1, 1, 'messsage two meta public channel', false),
(1, 2, 'message one story private channel', true),
(1, 2, 'messsage two meta private channel', false);
