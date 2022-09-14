package storage

type Event struct {
	ID          int64
	Title       string
	DateStart   string
	DateEnd     string
	Description string
	UserId      int
	WhenNotify  uint
}
