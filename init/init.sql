CREATE TABLE IF NOT EXISTS banners (
    id SERIAL PRIMARY KEY,
    tag_ids INTEGER[],
    feature_id INTEGER,
    content JSONB,
    is_active BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tag_ids ON banners USING GIN (tag_ids);
