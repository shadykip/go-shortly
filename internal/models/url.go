package models

type URL struct {
	ID        uint   `gorm:"primaryKey"`
	Original  string `gorm:"not null;uniqueIndex"`
	ShortCode string `gorm:"not null;uniqueIndex;size:8"`
	Clicks    int64  `gorm:"default:0"`
}
