CREATE TABLE IF NOT EXISTS banners (
    id SERIAL PRIMARY KEY,
    feature_id INT UNIQUE NOT NULL,
    content TEXT NOT NULL,
    is_active BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (feature_id) REFERENCES features(id) ON DELETE CASCADE
);