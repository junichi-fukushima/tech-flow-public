-- User tables
CREATE TABLE users
(
    id            VARCHAR(255) NOT NULL PRIMARY KEY,
    session_token VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) DEFAULT CHARSET = utf8mb4;

-- Feed tables
CREATE TABLE feeds
(
    id              BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    link            TEXT         NOT NULL,
    description     TEXT         NULL,
    category        TEXT         NULL,
    image           TEXT         NULL,
    language        VARCHAR(10)  NULL,
    last_build_date DATETIME     NULL,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE categories
(
    id   BIGINT       NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE tags
(
    id          BIGINT       NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(255) NOT NULL UNIQUE,
    category_id BIGINT       NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories (id)
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE articles
(
    id          BIGINT       NOT NULL PRIMARY KEY AUTO_INCREMENT,
    feed_id     BIGINT       NOT NULL,
    category_id BIGINT       NOT NULL,
    title       VARCHAR(255) NOT NULL,
    link        TEXT         NOT NULL,
    description TEXT         NULL,
    pub_date    DATETIME     NOT NULL,
    guid        VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE article_tags
(
    id         BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    article_id BIGINT NOT NULL,
    tag_id     BIGINT NOT NULL,
    FOREIGN KEY (article_id) REFERENCES articles (id),
    FOREIGN KEY (tag_id) REFERENCES tags (id)
) DEFAULT CHARSET = utf8mb4;


-- Event tables
CREATE TABLE item_metadata_events
(
    id         VARCHAR(255) NOT NULL PRIMARY KEY,
    timestamp  TIMESTAMP    NOT NULL,
    fields     JSON         NULL,
    article_id BIGINT       NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES articles (id)
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE user_metadata_events
(
    id         VARCHAR(255) NOT NULL PRIMARY KEY,
    timestamp  TIMESTAMP    NOT NULL,
    fields     JSON         NULL,
    user_id    VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE ranking_events
(
    id         VARCHAR(255) NOT NULL PRIMARY KEY,
    timestamp  TIMESTAMP    NOT NULL,
    fields     JSON         NULL,
    user_id    VARCHAR(255) NOT NULL,
    articles   JSON         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
) DEFAULT CHARSET = utf8mb4;

CREATE TABLE interaction_events
(
    id                     VARCHAR(255) NOT NULL PRIMARY KEY,
    timestamp              TIMESTAMP    NOT NULL,
    fields                 JSON         NULL,
    user_id                VARCHAR(255) NOT NULL,
    ranking_event_id       VARCHAR(255) NULL,
    item_metadata_event_id VARCHAR(255) NOT NULL,
    event_type             VARCHAR(255) NOT NULL,
    created_at             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (ranking_event_id) REFERENCES ranking_events (id),
    FOREIGN KEY (item_metadata_event_id) REFERENCES item_metadata_events (id)
) DEFAULT CHARSET = utf8mb4;
