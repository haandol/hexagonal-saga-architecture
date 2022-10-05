-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hotel_bookings (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  trip_id BIGINT NOT NULL,
  hotel_id BIGINT NOT NULL,
  status VARCHAR(32) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX hotel_bookings_trip_id ON hotel_bookings (trip_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX hotel_bookings_id_status ON hotel_bookings (id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE hotel_bookings
-- +goose StatementEnd
