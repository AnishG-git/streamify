-- +goose Up
-- +goose StatementBegin
CREATE TABLE room (
    code CHAR(5) PRIMARY KEY,
    participants INT NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE room;
-- +goose StatementEnd
