-- +goose Up
-- +goose StatementBegin
CREATE TABLE sagas (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  correlation_id VARCHAR(36) NOT NULL UNIQUE,
  trip_id BIGINT NOT NULL,
  car_id BIGINT NOT NULL,
  hotel_id BIGINT NOT NULL,
  flight_id BIGINT NOT NULL,
  status VARCHAR(16) NOT NULL,
  history JSON NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_created_at ON sagas (created_at);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_corr_id ON sagas (correlation_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sagas
-- +goose StatementEnd
