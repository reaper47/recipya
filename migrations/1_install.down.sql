DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO recipya;
GRANT ALL ON SCHEMA public TO public;

CREATE TABLE public.schema_migrations (
	"version" int8 NOT NULL,
	dirty bool NOT NULL,
	CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);
