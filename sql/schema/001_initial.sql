-- +goose Up
-- +goose StatementBegin
CREATE TABLE "application" (
    "id" UUID DEFAULT gen_random_uuid(),
    "name" VARCHAR(64) NOT NULL UNIQUE,
    PRIMARY KEY("id")
);

INSERT INTO "application"("name") VALUES('api_demo');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "application";
-- +goose StatementEnd