package entities

type Debtor struct {
	Username string
}

type NotificationMessage struct {
	ChatID    int64
	Username  string
	Text      string
	StickerID string
}
