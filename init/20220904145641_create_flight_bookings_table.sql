-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS flight_bookings (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  trip_id BIGINT NOT NULL,
  flight_id BIGINT NOT NULL,
  status VARCHAR(32) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX flight_bookings_trip_id ON flight_bookings (trip_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX flight_bookings_id_status ON flight_bookings (id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE flight_bookings
-- +goose StatementEnd
