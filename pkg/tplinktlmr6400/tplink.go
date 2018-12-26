package tplinktlmr6400

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/joonas-fi/modemrebooter/pkg/mrtypes"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	errSessionHashNotFound         = errors.New("session hash not found")
	errRestartConfirmationNotFound = errors.New("restart confirmation not found")
)

var httpClient = &http.Client{}

func New() mrtypes.ModemRebooter { return &tplink{} }

type tplink struct{}

func (r *tplink) Reboot(_ context.Context, conf mrtypes.Config) error {
	sessionHash, err := tryLogin(conf)
	if err != nil {
		// login fails sometimes, but works on 2nd try (because why not)
		sessionHash, err = tryLogin(conf)
		if err != nil {
			return err
		}
	}

	// now make a request to the reboot endpoint
	req, err := http.NewRequest("GET", rebootEndpoint(sessionHash, conf), nil)
	if err != nil {
		return err
	}

	// trololollol the reboot endpoint uses Referer as a security check
	req.Header.Set("Referer", homepageEndpoint(sessionHash, conf))

	attachAuthorizationCookie(req, conf)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !findRestartConfirmationFromResponse(string(body)) {
		return errRestartConfirmationNotFound
	}

	return nil
}

func tryLogin(conf mrtypes.Config) (string, error) {
	req, err := http.NewRequest("GET", loginEndpoint(conf), nil)
	if err != nil {
		return "", err
	}

	attachAuthorizationCookie(req, conf)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return findSessionHashFromLoginResponse(string(body))
}

var findSessionHashFromLoginResponseRe = regexp.MustCompile("([A-Z]+)\\/userRpm\\/Index.htm")

func findSessionHashFromLoginResponse(body string) (string, error) {
	matches := findSessionHashFromLoginResponseRe.FindStringSubmatch(body)
	if matches == nil {
		return "", errSessionHashNotFound
	}

	return matches[1], nil
}

func findRestartConfirmationFromResponse(body string) bool {
	return strings.Contains(body, "Restarting...") && strings.Contains(body, "Please wait a moment")
}

func fuckUpAuthorizationCookie(username string, password string) string {
	// from their JavaScript:
	//     password = hex_md5($("pcPassword").value);
	//     var auth = "Basic "+ Base64Encoding(username + ":" + password);
	//     document.cookie = "Authorization="+escape(auth)+";path=/";

	userAndPassword := fmt.Sprintf("%s:%x", username, md5.Sum([]byte(password)))

	authString := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(userAndPassword)))

	return strings.Replace(url.QueryEscape(authString), "+", "%20", -1)
}

func attachAuthorizationCookie(req *http.Request, conf mrtypes.Config) {
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: fuckUpAuthorizationCookie(conf.AdminUsername, conf.AdminPassword),
	})
}

func loginEndpoint(conf mrtypes.Config) string {
	return conf.Address + "/userRpm/LoginRpm.htm?Save=Save"
}

// this weird fucker has session ID in its URL (they probably think it's secure because
// security implementation details are in URL LOLOLOL). welcome to 90s, baby
func homepageEndpoint(sessionHash string, conf mrtypes.Config) string {
	return fmt.Sprintf("%s/%s/userRpm/Index.htm", conf.Address, sessionHash)
}

func rebootEndpoint(sessionHash string, conf mrtypes.Config) string {
	return fmt.Sprintf("%s/%s/userRpm/SysRebootRpm.htm?Reboot=Reboot", conf.Address, sessionHash)
}
