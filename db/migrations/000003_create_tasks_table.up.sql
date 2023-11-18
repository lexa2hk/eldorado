CREATE TABLE IF NOT EXISTS "public".tasks (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    title varchar(100) NOT NULL,
    description varchar(255) NOT NULL,
    is_completed boolean DEFAULT false NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    created_on timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_on timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_on timestamp,
    CONSTRAINT pk_tasks PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_tasks_user ON "public".tasks (user_id);
ALTER TABLE "public".tasks
ADD CONSTRAINT fk_tasks_users FOREIGN KEY (user_id) REFERENCES "public".users(id);
