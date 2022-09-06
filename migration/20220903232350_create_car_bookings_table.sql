-- +goose Up
-- +goose StatementBegin
CREATE TABLE car_bookings (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  trip_id BIGINT NOT NULL,
  car_id BIGINT NOT NULL,
  status VARCHAR(16) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_trip_id ON car_bookings (trip_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_id_status ON car_bookings (id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE car_bookings
-- +goose StatementEnd
