-- +goose Up
-- +goose StatementBegin
CREATE TABLE hotel_bookings (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  trip_id BIGINT NOT NULL,
  hotel_id BIGINT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_trip_id ON hotel_bookings (trip_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_hotel_id ON hotel_bookings (hotel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE hotel_bookings
-- +goose StatementEnd
