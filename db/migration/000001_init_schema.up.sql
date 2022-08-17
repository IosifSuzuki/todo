CREATE TABLE account
(
    id            serial PRIMARY KEY,
    username      VARCHAR(50) UNIQUE NOT NULL,
    hash_password TEXT               NOT NULL,
    email         TEXT               NOT NULL,
    created_on    TIMESTAMP          NOT NULL DEFAULT current_timestamp
);

CREATE TABLE item
(
    id          serial PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    created_on  TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    updated_on  TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    closed      BOOLEAN      NOT NULL
);

CREATE TABLE account_item
(
    account_id INT NOT NULL,
    item_id    INT NOT NULL,
    CONSTRAINT account_fk FOREIGN KEY (account_id) REFERENCES account (id),
    CONSTRAINT item_fk FOREIGN KEY (item_id) REFERENCES item (id)
);

CREATE INDEX ON account (username);
CREATE INDEX ON item (id);