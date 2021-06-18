CREATE TABLE IF NOT EXISTS postgres
(
    key text unique NOT NULL,
    value text  NOT NULL,
    PRIMARY KEY (key)
);

