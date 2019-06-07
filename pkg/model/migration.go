package model

// Migration defines model for migrations
type Migration struct {
	ID      int64
	Version int64
	Changes string
}
