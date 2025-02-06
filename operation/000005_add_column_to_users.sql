ALTER TABLE users
ADD COLUMN has_favorite_categories BOOLEAN DEFAULT FALSE NOT NULL
AFTER session_token;
