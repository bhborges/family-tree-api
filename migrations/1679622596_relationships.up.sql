CREATE TABLE IF NOT EXISTS "relationships" (
	"id" uuid NOT NULL DEFAULT uuid_generate_v4(),
	"parent_id" uuid,
	"children_id" uuid,
	"created_at" timestamptz NOT NULL DEFAULT NOW(),
	"updated_at" timestamptz NOT NULL DEFAULT NOW(),
	"deleted_at" timestamptz,
	PRIMARY KEY ("id")
);