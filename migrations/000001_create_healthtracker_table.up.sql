CREATE TABLE IF NOT EXISTS healthtracker (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
calories integer NOT NULL,
walking integer NOT NULL,
hydrate float NOT NULL,
sleep integer NOT NULL, 
user_id integer NOT NULL, 
version integer NOT NULL DEFAULT 1
);