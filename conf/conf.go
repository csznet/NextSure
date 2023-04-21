package conf

type Link struct {
	Lid   uint `gorm:"primaryKey;autoIncrement"`
	Title string
	Url   string
	Img   string
	Date  int64 `gorm:"not null"`
}
