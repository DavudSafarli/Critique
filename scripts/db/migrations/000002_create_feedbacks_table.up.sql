CREATE TABLE IF NOT EXISTS feedbacks (
    id serial PRIMARY KEY,
    title VARCHAR (255) NOT NULL,
    body VARCHAR (65535) NOT NULL,
    created_by VARCHAR(255),
    created_at timestamptz NOT NULL
);