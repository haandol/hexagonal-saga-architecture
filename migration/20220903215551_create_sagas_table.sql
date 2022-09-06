-- +goose Up
-- +goose StatementBegin
CREATE TABLE sagas (
  id SERIAL PRIMARY KEY,
  correlation_id VARCHAR(36) NOT NULL UNIQUE,
  trip_id BIGINT NOT NULL,
  car_id BIGINT NOT NULL,
  hotel_id BIGINT NOT NULL,
  flight_id BIGINT NOT NULL,
  car_booking_id BIGINT NOT NULL,
  hotel_booking_id BIGINT NOT NULL,
  flight_booking_id BIGINT NOT NULL,
  status VARCHAR(32) NOT NULL,
  history JSON NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX sagas_trip_id ON sagas (trip_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX sagas_corr_id ON sagas (correlation_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX sagas_corr_id_status ON sagas (correlation_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sagas
-- +goose StatementEnd
