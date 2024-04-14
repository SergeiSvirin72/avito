CREATE TABLE IF NOT EXISTS banner_tag (
    banner_id INT NOT NULL,
    tag_id INT NOT NULL,
    FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (banner_id, tag_id)
);