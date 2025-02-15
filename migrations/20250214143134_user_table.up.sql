CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    hashed_token VARCHAR(500) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 day'
);

CREATE TABLE report_jobs (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    report_type VARCHAR NOT NULL,
    output_file_path VARCHAR,
    download_url VARCHAR,
    download_url_expires_at TIMESTAMPTZ,
    error_message VARCHAR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    failed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    PRIMARY KEY(user_id, id)
);