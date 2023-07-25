BEGIN TRANSACTION;
CREATE TYPE processing_state as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TYPE accrual_state as enum ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED');
CREATE TABLE IF NOT EXISTS accumulations
(
    id                serial primary key,
    user_id           integer                  not null,
    order_number      decimal                  null,
    uploaded_at       timestamp with time zone not null,
    processing_status processing_state         null,
    accrual_status    accrual_state            null,
    amount            float8                   null,
    processed_at      timestamp with time zone null
);
CREATE INDEX accumulations_user_id_index on accumulations (user_id);
COMMIT;