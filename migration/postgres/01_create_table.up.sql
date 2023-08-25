CREATE TABLE orders (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR,
    "phone" VARCHAR,
    "lat" VARCHAR,
    "long" VARCHAR,
    "address" VARCHAR,
    "photo" VARCHAR,
    "created_at" TIMESTAMP default CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);