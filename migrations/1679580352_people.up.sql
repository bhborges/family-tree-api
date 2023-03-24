CREATE TABLE IF NOT EXISTS "people" (
	"id" uuid NOT NULL DEFAULT uuid_generate_v4(),
	"name" varchar(255) NOT NULL,
	"created_at" timestamptz NOT NULL DEFAULT NOW(),
	"updated_at" timestamptz NOT NULL DEFAULT NOW(),
	"deleted_at" timestamptz,
	PRIMARY KEY ("id")
);
