CREATE SCHEMA IF NOT EXISTS monitoring;

SET search_path TO monitoring, public;


CREATE TABLE IF NOT EXISTS Users(
	id int NOT NULL,
	"username" text NOT NULL,
	"password" text NOT NULL,
	"role" int8 NOT NULL
);

CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    address VARCHAR(255),
    method VARCHAR(255),
    header JSONB, -- Assuming header is a JSON object
    body JSONB,   -- Assuming body is a JSON object
    access_level INTEGER,
    execution_time INTEGER,
    allowed_users TEXT[]
);
