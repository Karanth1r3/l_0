BEGIN;

CREATE TABLE IF NOT EXISTS service_storage 
(
    id bigserial NOT NULL
        CONTRAINT service_storage_pk
            PRIMARY KEY,
    order_uid VARCHAR NOT NULL,
    value JSONB NOT NULL
);

COMMENT ON TABLE service_storage IS 'KV-storage for records';

COMMENT ON COLUMN service_storage.id IS 'Autoincrement id';

COMMENT ON COLUMN service_storage.order_uid 'A unitque key in the table';

COMMENT ON COLUMN sevice_storage.value IS 'Static data';

CREATE UNIQUE INDEX IF NOT EXISTS service_storage_key_uindex
    ON service_storage (order_uid);

COMMIT;