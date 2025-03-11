CREATE TABLE messages (
	id serial NOT NULL,
    user_id varchar(10),
	session_id varchar(32) NOT NULL,
    message_id uuid NULL,
	query text NOT NULL,
	answer text NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
	is_safe bool DEFAULT false NOT NULL,
    is_like bool DEFAULT false NOT NULL,
	CONSTRAINT message_pkey PRIMARY KEY (id)
);