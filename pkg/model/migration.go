package model

// Migration defines model for migrations
type Migration struct {
	ID      int64
	Version string
	Changes string
}
