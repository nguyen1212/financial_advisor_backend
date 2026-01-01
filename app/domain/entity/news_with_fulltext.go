package entity

type NewsWithFullText struct {
	ID      uint64 `gorm:"primaryKey"`
	NewsID  string `gorm:"column:news_id"`
	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`
}

func (NewsWithFullText) TableName() string {
	return "news"
}
