-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS trips (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT NOT NULL,
  car_id BIGINT NOT NULL,
  hotel_id BIGINT NOT NULL,
  flight_id BIGINT NOT NULL,
  car_booking_id BIGINT NOT NULL,
  hotel_booking_id BIGINT NOT NULL,
  flight_booking_id BIGINT NOT NULL,
  status VARCHAR(32) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX trips_created_at ON trips (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trips
-- +goose StatementEnd
