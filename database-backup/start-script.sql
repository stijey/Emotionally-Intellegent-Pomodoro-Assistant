-- psql -U thefriyia -a -f database-backup/start-script.sql

-- DROP SCHEMA temp CASCADE;
-- CREATE SCHEMA temp;
-- CREATE TABLE temp.userinfo
-- (
--     uid serial NOT NULL,
--     username character varying(100) NOT NULL,
--     departname character varying(500) NOT NULL,
--     Created date,
--     CONSTRAINT userinfo_pkey PRIMARY KEY (uid)
-- )
-- WITH (OIDS=FALSE);

DROP SCHEMA pomodoro CASCADE;
DROP DATABASE pomodoro;
DROP TABLE users;
CREATE DATABASE pomodoro;
CREATE TABLE users
(
    uid SERIAL,
    username text NOT NULL UNIQUE,
    password bytea NOT NULL,
    weekly_goals bytea,
    CONSTRAINT users_pkey PRIMARY KEY (uid)
) WITH (OIDS=FALSE);

INSERT INTO users (username,password)
VALUES ('dan','pass');