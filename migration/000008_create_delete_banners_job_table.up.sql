CREATE TABLE IF NOT EXISTS delete_banners_job (
    id SERIAL PRIMARY KEY,
    feature_id INT NULL,
    tag_id INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (feature_id) REFERENCES features(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);