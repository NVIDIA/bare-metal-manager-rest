DROP INDEX IF EXISTS component_machine_id;

ALTER TABLE component
    DROP COLUMN IF EXISTS machine_id,
    DROP COLUMN IF EXISTS power_state;
