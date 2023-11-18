CREATE TABLE IF NOT EXISTS "public".users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    email varchar(150) NOT NULL,
    username varchar(150) NOT NULL,
    name varchar(150) NOT NULL,
    encrypted_password varchar(1024) NOT NULL,
    created_on timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_on timestamp,
    CONSTRAINT pk_users PRIMARY KEY (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS unq_users_email ON "public".users (email);
CREATE UNIQUE INDEX IF NOT EXISTS unq_users_username ON "public".users (username);
