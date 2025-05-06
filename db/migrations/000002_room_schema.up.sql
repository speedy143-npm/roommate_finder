CREATE TABLE rooms (
    "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
    "user_id" VARCHAR(36) REFERENCES "users"("id") ON DELETE CASCADE, -- Owner of listing
    "location" VARCHAR(255) NOT NULL,
    "rent_price" DECIMAL(10,2) NOT NULL,
    "room_details" JSONB, -- Stores amenities, size, etc.
    "availability_status" BOOLEAN DEFAULT TRUE,
    "created_at" TIMESTAMP DEFAULT now()
);