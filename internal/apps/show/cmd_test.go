package show

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestShowMax(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	if err := run(); err != nil {
		t.Fatal(err)
	}
}
