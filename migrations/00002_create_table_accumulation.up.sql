BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS accumulation
(
    id                   serial primary key,
    user_id              integer   not null,
    order_number         integer   null,
    date                 timestamp not null,
    processing_status_id integer   not null,
    accrual_status_id    integer   not null,

);
CREATE TABLE IF NOT EXISTS processing_statuses
(
    id   serial primary key,
    name varchar(30) not null
);
CREATE TABLE IF NOT EXISTS accrual_statuses
(
    id   serial primary key,
    name varchar(30) not null
);
COMMIT;