package entity

import (
	"time"

	"github.com/haandol/hexagonal/internal/dto"
	"gorm.io/datatypes"
)

type Outbox struct {
	ID         uint           `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	KafkaTopic string         `gorm:"type:varchar(256);<-:create;"`
	KafkaKey   string         `gorm:"type:varchar(100);<-:create;"`
	KafkaValue datatypes.JSON `gorm:"type:json;<-:create;"`
	IsSent     bool           `gorm:"type:bool;default:false;"`
	CreatedAt  time.Time      `gorm:"type:timestamp;<-:create;"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;"`
}

type Outboxes []*Outbox

func (m *Outbox) DTO() dto.Outbox {
	return dto.Outbox{
		ID:         m.ID,
		KafkaTopic: m.KafkaTopic,
		KafkaKey:   m.KafkaKey,
		KafkaValue: string(m.KafkaValue),
		IsSent:     m.IsSent,
		CreatedAt:  m.CreatedAt,
	}
}

func (m Outboxes) DTO() []dto.Outbox {
	outboxes := make([]dto.Outbox, 0, len(m))
	for _, outbox := range m {
		outboxes = append(outboxes, outbox.DTO())
	}
	return outboxes
}
