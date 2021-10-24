CREATE TABLE public.users (
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	deleted_at timestamptz NULL,
	name varchar NOT NULL,
	"data" varchar NULL,
	perms int2 NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);