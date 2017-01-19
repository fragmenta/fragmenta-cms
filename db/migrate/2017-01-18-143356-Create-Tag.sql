DROP TABLE IF EXISTS tags;
CREATE TABLE tags (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
parent_id integer,
name text,
summary text,
url text,
sort integer,
dotted_ids text,
status integer
);
ALTER TABLE tags OWNER TO fragmentaweb_server;
