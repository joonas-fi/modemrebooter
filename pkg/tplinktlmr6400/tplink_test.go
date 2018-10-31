package tplinktlmr6400

import (
	"github.com/function61/gokit/assert"
	"github.com/joonas-fi/modemrebooter/pkg/mrtypes"
	"testing"
)

func TestFuckUpAuthorizationCookie(t *testing.T) {
	assert.EqualString(t, fuckUpAuthorizationCookie("admin", "1234"), "Basic%20YWRtaW46ODFkYzliZGI1MmQwNGRjMjAwMzZkYmQ4MzEzZWQwNTU%3D")
}

func TestFindSessionHashFromLoginResponse(t *testing.T) {
	hash, err := findSessionHashFromLoginResponse(`<body><script language="javaScript">window.parent.location.href = "http://192.168.1.1/ISAQWMDBPYGLRYQC/userRpm/Index.htm";</script></body></html>`)
	assert.True(t, err == nil)
	assert.EqualString(t, hash, "ISAQWMDBPYGLRYQC")
}

func TestEndpoints(t *testing.T) {
	cfg := mrtypes.Config{
		Type:          "tplinktlmr6400",
		Address:       "http://192.168.1.1",
		AdminUsername: "admin",
		AdminPassword: "1234",
	}

	assert.EqualString(t, loginEndpoint(cfg), "http://192.168.1.1/userRpm/LoginRpm.htm?Save=Save")

	assert.EqualString(t, rebootEndpoint("ISAQWMDBPYGLRYQC", cfg), "http://192.168.1.1/ISAQWMDBPYGLRYQC/userRpm/SysRebootRpm.htm?Reboot=Reboot")
}

func TestFindRestartConfirmationFromResponse(t *testing.T) {
	correctBody := "<TR>\n<TD class=h2 id=\"t_restart\">Restarting...</TD>\n</TR>\n<TR><TD class = \"h2\" id =\"t_notice\" style=\"display:none\">Please wait a moment, if the browser ..."
	assert.True(t, findRestartConfirmationFromResponse(correctBody) == true)

	incorrectBody := "<TR>\n<TD class=h2 id=\"t_restart\">Restarting...</TD>\n</TR>\n<TR><TD class = \"h2\" id =\"t_notice\" style=\"display:none\">Please wait a m0ment, if the browser ..."
	assert.True(t, findRestartConfirmationFromResponse(incorrectBody) == false)
}
