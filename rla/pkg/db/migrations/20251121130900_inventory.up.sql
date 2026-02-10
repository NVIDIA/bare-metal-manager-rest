ALTER TABLE component 
    ADD COLUMN machine_id TEXT,
    ADD COLUMN power_state INT;

CREATE INDEX component_machine_id ON component (machine_id);
