DROP TABLE IF EXISTS pages;
CREATE TABLE pages (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
status integer,
author_id integer,
url text,
name text,
summary text,
keywords text,
template text,
text text
);
ALTER TABLE pages OWNER TO fragmentaweb_server;
