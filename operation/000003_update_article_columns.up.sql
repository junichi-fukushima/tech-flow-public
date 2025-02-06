ALTER TABLE articles MODIFY COLUMN guid VARCHAR(512) NOT NULL UNIQUE;
ALTER TABLE articles ADD UNIQUE INDEX idx_guid (guid);
ALTER TABLE articles ADD INDEX idx_feed_category (feed_id, category_id);