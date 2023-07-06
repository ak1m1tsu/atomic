CREATE TABLE IF NOT EXISTS alias (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name varchar(50) NOT NULL,
    url text NOT NULL,
    CONSTRAINT pk_alias PRIMARY KEY (id),
    CONSTRAINT unq_alias_name UNIQUE (name)
)
