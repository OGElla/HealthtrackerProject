CREATE TABLE IF NOT EXISTS goals (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
walking integer NOT NULL,
achieved bool NOT NULL DEFAULT false,
user_id integer NOT NULL, 
version integer NOT NULL DEFAULT 1
);