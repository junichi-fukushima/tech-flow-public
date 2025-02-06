-- NOTE: 新規作成済の場合流す必要はない！
-- Drop FK
ALTER TABLE ranking_events DROP FOREIGN KEY fk_ranking_events;
-- Rename user_id
ALTER TABLE ranking_events RENAME COLUMN user_metadata_event_id TO user_id;
-- Add new FK
ALTER TABLE ranking_events ADD CONSTRAINT ranking_events_ibfk_1 FOREIGN KEY (user_id) REFERENCES users(id);

-- Drop FK
ALTER TABLE interaction_events DROP FOREIGN KEY fk_interaction_events;
-- Rename user_id
ALTER TABLE interaction_events RENAME COLUMN user_metadata_event_id TO user_id;
-- Add new FK
ALTER TABLE interaction_events ADD CONSTRAINT interaction_events_ibfk_1 FOREIGN KEY (user_id) REFERENCES users(id);
