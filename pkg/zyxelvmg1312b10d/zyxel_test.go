package zyxelvmg1312b10d

import (
	"github.com/function61/gokit/assert"
	"testing"
)

func TestAuthCookieString(t *testing.T) {
	assert.EqualString(t, authCookieString("admin", "12345"), "YWRtaW46MTIzNDU%3D")
}
