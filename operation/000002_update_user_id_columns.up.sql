-- Drop FK
ALTER TABLE ranking_events DROP FOREIGN KEY ranking_events_ibfk_1;
-- Rename user_id
ALTER TABLE ranking_events RENAME COLUMN user_id TO user_metadata_event_id;
-- Add new FK
ALTER TABLE ranking_events ADD CONSTRAINT fk_ranking_events FOREIGN KEY (user_metadata_event_id) REFERENCES user_metadata_events(id);

-- Drop FK
ALTER TABLE interaction_events DROP FOREIGN KEY interaction_events_ibfk_1;
-- Rename user_id
ALTER TABLE interaction_events RENAME COLUMN user_id TO user_metadata_event_id;
-- Add new FK
ALTER TABLE interaction_events ADD CONSTRAINT fk_interaction_events FOREIGN KEY (user_metadata_event_id) REFERENCES user_metadata_events(id);
