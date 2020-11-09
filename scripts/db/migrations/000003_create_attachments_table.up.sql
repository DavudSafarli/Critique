CREATE TABLE IF NOT EXISTS attachments (
    id serial PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    path VARCHAR (65535) NOT NULL,
    feedback_id INT NOT NULL,
    CONSTRAINT fk_feedback
        FOREIGN KEY(feedback_id)
            REFERENCES feedbacks(id)
            ON DELETE CASCADE
);