BEGIN TRANSACTION;
CREATE TYPE processing_state as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TYPE accrual_state as enum ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED');
CREATE TABLE IF NOT EXISTS accumulation
(
    id                serial primary key,
    user_id           integer                  not null,
    order_number      integer                  null,
    uploaded_at       timestamp with time zone not null,
    processing_status processing_state         not null,
    accrual_status    accrual_state            null,
    amount            integer                  null,
    processed_at      timestamp with time zone null
);
CREATE INDEX accumulation_user_id_index on accumulation (user_id);
COMMIT;