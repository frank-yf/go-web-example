package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInternetAddress(t *testing.T) {
	ip, err := GetInternetAddress()
	assert.Nil(t, err)
	t.Log(ip)
}
