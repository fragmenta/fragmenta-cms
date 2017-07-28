DROP TABLE IF EXISTS redirects;
CREATE TABLE redirects (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
new_url text,
old_url text
);
ALTER TABLE redirects OWNER TO fragmentaweb_server;
