-- +goose Up
-- +goose StatementBegin
CREATE TABLE car_bookings (
  id SERIAL PRIMARY KEY,
  trip_id BIGINT NOT NULL,
  car_id BIGINT NOT NULL,
  status VARCHAR(16) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX car_bookings_trip_id ON car_bookings (trip_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX car_bookings_id_status ON car_bookings (id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE car_bookings
-- +goose StatementEnd
