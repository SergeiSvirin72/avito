CREATE TABLE IF NOT EXISTS banner_history (
    id SERIAL PRIMARY KEY,
    banner_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE
);