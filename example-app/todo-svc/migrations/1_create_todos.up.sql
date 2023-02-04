CREATE TABLE todos
(
    id        SERIAL PRIMARY KEY,
    user_id   INTEGER                     NOT NULL,
    title     TEXT                        NOT NULL,
    created   TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    completed BOOLEAN                     NOT NULL DEFAULT FALSE
);

CREATE INDEX todos_user_id_idx ON todos (user_id, id);
