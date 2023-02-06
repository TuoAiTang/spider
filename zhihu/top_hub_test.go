package zhihu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuoaitang/spider/log"
)

func TestGetTopHub(t *testing.T) {
	hub, err := GetTopHub()
	assert.Nil(t, err)
	log.Info("hub:%v", hub.String())
}
