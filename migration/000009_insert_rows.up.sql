DELETE FROM delete_banners_job;
DELETE FROM banner_history;
DELETE FROM users;
DELETE FROM banner_tag;
DELETE FROM banners;
DELETE FROM features;
DELETE FROM tags;

INSERT INTO tags (id) SELECT * FROM generate_series(1, 1000);
SELECT setval('tags_id_seq', 1000);

INSERT INTO features (id) SELECT * FROM generate_series(1, 1000);
SELECT setval('features_id_seq', 1000);


INSERT INTO banners (id, feature_id, content, is_active) VALUES
    (1, 1, '{"title_1": "title_1", "text_1": "text_1", "url_1": "url_1"}', true),
    (2, 2, '{"title_2": "title_2", "text_2": "text_2", "url_2": "url_2"}', true),
    (3, 3, '{"title_3": "title_3", "text_3": "text_3", "url_3": "url_3"}', true),
    (4, 4, '{"title_4": "title_4", "text_4": "text_4", "url_4": "url_4"}', false),
    (5, 5, '{"title_5": "title_5", "text_5": "text_5", "url_5": "url_5"}', false),
    (6, 6, '{"title_6": "title_6", "text_6": "text_6", "url_6": "url_6"}', false)
;
SELECT setval('banners_id_seq', 6);

INSERT INTO banner_tag (banner_id, tag_id) VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (2, 1),
    (2, 2),
    (3, 3),
    (4, 1),
    (4, 2),
    (4, 3),
    (5, 1),
    (5, 2),
    (6, 3)
;

INSERT INTO users (email, password, role) VALUES
    ('admin@test.com', '$2a$04$0N577rXqpu1tA5uIt/X7VeF2bCgoAYF2OT.kvv67jgFxgPy9DK0hi', 'admin'),
    ('user1@test.com', '$2a$04$0N577rXqpu1tA5uIt/X7VeF2bCgoAYF2OT.kvv67jgFxgPy9DK0hi', 'user'),
    ('user2@test.com', '$2a$04$0N577rXqpu1tA5uIt/X7VeF2bCgoAYF2OT.kvv67jgFxgPy9DK0hi', 'user'),
    ('user3@test.com', '$2a$04$0N577rXqpu1tA5uIt/X7VeF2bCgoAYF2OT.kvv67jgFxgPy9DK0hi', 'user'),
    ('user4@test.com', '$2a$04$0N577rXqpu1tA5uIt/X7VeF2bCgoAYF2OT.kvv67jgFxgPy9DK0hi', 'user'),
    ('user5@test.com', '$2a$04$0N577rXqpu1tA5uIt/X7VeF2bCgoAYF2OT.kvv67jgFxgPy9DK0hi', 'user')
;

INSERT INTO banner_history (banner_id, content) VALUES
    (1, '{"title_1": "version_1", "text_1": "text_1", "url_1": "url_1"}'),
    (1, '{"title_1": "version_2", "text_1": "text_1", "url_1": "url_1"}'),
    (1, '{"title_1": "version_3", "text_1": "text_1", "url_1": "url_1"}'),
    (2, '{"title_2": "version_1", "text_2": "text_2", "url_2": "url_2"}'),
    (2, '{"title_2": "version_2", "text_2": "text_2", "url_2": "url_2"}'),
    (3, '{"title_3": "version_1", "text_3": "text_3", "url_3": "url_3"}')
;
SELECT setval('banners_id_seq', 6);
