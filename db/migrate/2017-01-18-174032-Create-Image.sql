DROP TABLE IF EXISTS images;
CREATE TABLE images (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
status integer,
author_id integer,
path text,
sort integer,
name text
);
ALTER TABLE images OWNER TO fragmentaweb_server;
