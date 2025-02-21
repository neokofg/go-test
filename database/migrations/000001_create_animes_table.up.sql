CREATE TABLE IF NOT EXISTS animes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    rating DECIMAL(3,1),
    episodes INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
