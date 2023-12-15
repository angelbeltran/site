BEGIN;


CREATE TABLE user_sessions (
	user_session_id 	SERIAL				PRIMARY KEY,
	user_id 					INTEGER 			NOT NULL REFERENCES users (user_id),
	expires_at				TIMESTAMPTZ 	NOT NULL DEFAULT NOW() + '1 hour'::INTERVAL
);
CREATE UNIQUE INDEX user_sessions_user_id_idx ON user_sessions (user_id);
CREATE INDEX user_sessions_expires_at_idx ON user_sessions (expires_at);


REASSIGN OWNED BY current_user TO pi;

COMMIT;
