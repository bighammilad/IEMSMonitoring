CREATE SCHEMA IF NOT EXISTS monitoring;

SET search_path TO monitoring, public;


CREATE TABLE IF NOT EXISTS Users(
	id int NOT NULL,
	"username" text NOT NULL,
	"password" text NOT NULL,
	"role" int8 NOT NULL
);
