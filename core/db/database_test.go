package db

import "testing"

func TestDatabase(t *testing.T) {
	Init("localhost", "postgres", "postgres", "postgres", 5433)
}
