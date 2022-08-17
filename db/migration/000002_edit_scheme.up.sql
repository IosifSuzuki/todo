ALTER TABLE account_item
DROP CONSTRAINT IF EXISTS account_fk;
ALTER TABLE account_item
    ADD CONSTRAINT account_fk
    FOREIGN KEY (account_id)
    REFERENCES account (id)
    ON DELETE CASCADE;
ALTER TABLE account_item
DROP CONSTRAINT IF EXISTS item_fk;
ALTER TABLE account_item
    ADD CONSTRAINT item_fk
    FOREIGN KEY (item_id)
    REFERENCES item (id)
    ON DELETE CASCADE;