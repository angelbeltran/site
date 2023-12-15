BEGIN;


CREATE TABLE users (
	user_id			SERIAL				PRIMARY KEY,
	name				TEXT					NOT NULL CONSTRAINT nonempty_user_name CHECK (name <> ''),
	password 		TEXT					NOT NULL CONSTRAINT nonempty_user_password CHECK (password <> ''),
	created_at	TIMESTAMPTZ		NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX users_name_idx ON users (name) INCLUDE (password);
CREATE INDEX users_created_at_idx ON users (created_at);


REASSIGN OWNED BY current_user TO pi;

COMMIT;
