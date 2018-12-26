package zyxelvmg1312b10d

import (
	"context"
	"encoding/base64"
	"github.com/function61/gokit/ezhttp"
	"github.com/joonas-fi/modemrebooter/pkg/mrtypes"
	"net/http"
	"net/url"
)

func New() mrtypes.ModemRebooter { return &zyxel{} }

type zyxel struct{}

func (r *zyxel) Reboot(conf mrtypes.Config) error {
	ctx, cancel := context.WithTimeout(context.TODO(), ezhttp.DefaultTimeout10s)
	defer cancel()

	_, err := ezhttp.Get(
		ctx,
		conf.Address+"/cgi-bin/Reboot",
		ezhttp.Cookie(http.Cookie{
			Name:  "Authentication",
			Value: authCookieString(conf.AdminUsername, conf.AdminPassword),
		}))
	return err
}

func authCookieString(username, password string) string {
	base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	// padding "=" shows up as "%3D"
	urlescaped := url.QueryEscape(base64Encoded)

	return urlescaped
}
