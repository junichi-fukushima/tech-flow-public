-- 一旦、既存レコードの存在を無視してinsertするためにIGNOREをつけてます
INSERT IGNORE INTO categories (id, name) VALUES
(1, 'プログラミング言語'),
(2, 'フレームワーク'),
(3, 'インフラ'),
(4, '最新技術'),
(5, 'AI'),
(6, 'IoT'),
(7, 'マネジメント'),
(8, 'NONE');