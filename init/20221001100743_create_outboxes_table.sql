-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outboxes (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  kafka_topic VARCHAR(256) NOT NULL,
  kafka_key VARCHAR(100) NOT NULL,
  kafka_value JSON,
  is_sent BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX outbox_is_sent ON outboxes (is_sent);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX outbox_kafka_topic_key ON outboxes (kafka_topic, kafka_key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outboxes
-- +goose StatementEnd
