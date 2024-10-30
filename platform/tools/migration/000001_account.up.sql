-- Up
CREATE TABLE gostarter_account
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(255)             NOT NULL,
    password   VARCHAR(255)             NOT NULL,
    email      VARCHAR(255) UNIQUE      NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Up
CREATE TABLE gostarter_role
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255)             NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Up
CREATE TABLE gostarter_account_role
(
    account_id INT,
    role_id    INT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (account_id, role_id),
    FOREIGN KEY (account_id) REFERENCES gostarter_account (id),
    FOREIGN KEY (role_id) REFERENCES gostarter_role (id)
);

-- Function to update timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for gostarter_account
CREATE TRIGGER update_gostarter_account_timestamp
    BEFORE UPDATE
    ON gostarter_account
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Trigger for gostarter_role
CREATE TRIGGER update_gostarter_role_timestamp
    BEFORE UPDATE
    ON gostarter_role
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();