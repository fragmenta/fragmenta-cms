DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
keywords text,
template text,
text text,
status integer,
author_id integer,
name text,
summary text
);
ALTER TABLE posts OWNER TO fragmentaweb_server;
