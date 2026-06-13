-- Create "jobs" table
CREATE TABLE "jobs" (
  "id" bigserial NOT NULL,
  "queue" integer NULL,
  "payload" text NULL,
  "status" text NULL,
  PRIMARY KEY ("id")
);
