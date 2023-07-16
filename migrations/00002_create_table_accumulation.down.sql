BEGIN TRANSACTION;
DROP TABLE IF EXISTS accumulation;
DROP TYPE IF EXISTS processing_state;
DROP TYPE IF EXISTS accrual_state;
DROP INDEX IF EXISTS accumulation_user_id_index;
COMMIT;