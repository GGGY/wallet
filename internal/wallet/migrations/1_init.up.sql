CREATE TYPE CURRENCY AS ENUM ('USD', 'EUR');

--account table
CREATE TABLE account
(
    id VARCHAR(256)                      NOT NULL
      CONSTRAINT account_pk PRIMARY KEY,
    balance NUMERIC(36, 18) DEFAULT 0    NOT NULL,
    currency CURRENCY                    NOT NULL,
    CONSTRAINT  balance_positive_check CHECK (balance >= 0)
);

--payment table
CREATE TABLE payment
(
    id            BIGSERIAL              NOT NULL
        CONSTRAINT payment_pk PRIMARY KEY,
    account VARCHAR(256)                 NOT NULL,
    amount NUMERIC(36, 18)    DEFAULT 0  NOT NULL,
    from_account VARCHAR(256) DEFAULT '',
    to_account VARCHAR(256)   DEFAULT '',
    CONSTRAINT payment_account_fk FOREIGN KEY (account)
    REFERENCES account(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);