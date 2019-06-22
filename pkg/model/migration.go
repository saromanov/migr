package model

// Migration defines model for migrations
type Migration struct {
	ID           int64
	Version      int64
	Changes      string
	Hash         *string
	Applied      bool
	ErrorMessage *string
	Failed       bool
	Status       string
}
