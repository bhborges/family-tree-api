CREATE TABLE IF NOT EXISTS "relationships" (
	"id" uuid NOT NULL DEFAULT uuid_generate_v4(),
	"parent_id" uuid,
	"child_id" uuid,
	PRIMARY KEY ("id"),
	FOREIGN KEY ("parent_id") REFERENCES "people" ("id"),
    FOREIGN KEY ("child_id") REFERENCES "people" ("id")
);
