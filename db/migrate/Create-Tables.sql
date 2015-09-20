/* Setup tables for cms */

CREATE TABLE fragmenta_metadata (
    id SERIAL NOT NULL,
    updated_at timestamp,
    fragmenta_version text,
    migration_version text,
    status int
);

ALTER TABLE fragmenta_metadata OWNER TO "[[.fragmenta_db_user]]";

DROP TABLE users;
CREATE TABLE users (
    id SERIAL NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    created_by int,
    updated_by int,
    status int,
    role int,
    email text,
    name text,
    title text,
    summary text,
    text text,
    image_id integer,
    encrypted_password text,
    reset_password_token text,
    password_reset_at timestamp,
    confirmation_token text,
    confirmation_token_at timestamp,
    confirmed_at timestamp
);
ALTER TABLE users OWNER TO "[[.fragmenta_db_user]]";


DROP TABLE pages;
CREATE TABLE pages (
    id SERIAL NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    created_by int,
    updated_by int,
    status int,
    url text,
    name text,
    summary text,
    keywords text, 
    text text default '<h1>Your Title</h1><p>Your text</p>',
    template text default 'pages/views/show.html.got'
);
ALTER TABLE pages OWNER TO "[[.fragmenta_db_user]]";


DROP TABLE tags;
CREATE TABLE tags (
    id SERIAL NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    created_by int,
    updated_by int,
    status int,
    parent_id int,
    url text,
    name text,
    summary text
);
ALTER TABLE tags OWNER TO "[[.fragmenta_db_user]]";

DROP TABLE pages_tags;
CREATE TABLE pages_tags (
   page_id int NOT NULL,
   tag_id int NOT NULL
);
ALTER TABLE pages_tags OWNER TO "[[.fragmenta_db_user]]";



DROP TABLE images;
CREATE TABLE images (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
status int,
name text,
path text,
sort integer
);
ALTER TABLE images OWNER TO "[[.fragmenta_db_user]]";

DROP TABLE images_pages;
CREATE TABLE images_pages (
   image_id int NOT NULL,
   page_id int NOT NULL
);
ALTER TABLE images_pages OWNER TO "[[.fragmenta_db_user]]";
