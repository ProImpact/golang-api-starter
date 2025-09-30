-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user" (
    "id" UUID DEFAULT gen_random_uuid(),
    "name" VARCHAR(64) NOT NULL UNIQUE,
    PRIMARY KEY("id")
);

INSERT INTO "user"("name") VALUES('admin');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
-- +goose StatementEnd