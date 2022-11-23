-- +goose Up
-- +goose StatementBegin
CREATE TABLE jobs (
    id              UUID       NOT NULL,
    created_at      TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP  NOT NULL,
    source_code_b64 TEXT       NOT NULL,
    status          VARCHAR(8) NOT NULL,
    exit_code       INT        NULL,

    PRIMARY KEY (id)
);
CREATE TABLE job_output_rows (
    id         SERIAL    NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    text       TEXT      NOT NULL,
    job_id     UUID      NOT NULL,

    PRIMARY KEY (id),

    FOREIGN KEY (job_id)
        REFERENCES jobs(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE job_output_rows;
DROP TABLE jobs;
-- +goose StatementEnd
