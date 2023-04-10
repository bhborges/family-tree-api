CREATE TABLE IF NOT EXISTS "people" (
	"id" uuid NOT NULL DEFAULT uuid_generate_v4(),
	"name" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);
