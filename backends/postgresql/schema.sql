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

-- Mapping table to determine which users are in what channels
CREATE TABLE channels_users (
    channelid bigserial references channels(id),
    userid bigserial references users(id),
    createdon timestamp default current_timestamp
);

CREATE TABLE messages (
    id bigserial primary key,
    userid bigserial references users(id),
    content varchar(200) NOT NULL,
    createdon timestamp default current_timestamp,
    lastupdated timestamp default current_timestamp
);

-- Use function provided from trigger.sql to handle updating lastmodified timestamps on updates
CREATE TRIGGER users_updated_at_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();
CREATE TRIGGER channels_updated_at_modtime BEFORE UPDATE ON channels FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();
CREATE TRIGGER messages_updated_at_modtime BEFORE UPDATE ON messages FOR EACH ROW EXECUTE PROCEDURE update_lastupdated_column();

-- Sample data
INSERT INTO users
(name, bio) VALUES
('user1', 'my first user');

INSERT INTO channels
(ownerid, dmid, name, description, isprivate) VALUES
(1, 1, 'my channel', 'my channel description', false);

INSERT INTO channels_users
(channelid, userid) VALUES
(1, 1)