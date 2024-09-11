package messaging

import (
	"github.com/lib/pq"
	"time"
)

type CreateProcess struct {
	WorkflowID string
	Name       string
	Desc       string
}

type Process struct {
	ID      string         `gorm:"unique;primaryKey;autoIncrement" json:"id"`
	UUID    string         `gorm:"index;type:text" json:"uuid"`
	Name    string         `gorm:"type:text" json:"name"`
	Desc    string         `gorm:"type:text" json:"desc"`
	Stage   string         `gorm:"type:text" json:"stage"`
	Cost    float32        `gorm:"type:float" json:"cost"`
	Parts   pq.StringArray `gorm:"type:text[]" json:"parts"`
	Created time.Time      `gorm:"type:timestamp" json:"created"`
}
