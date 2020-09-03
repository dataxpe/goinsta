package goinsta

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Instagram represent the main API handler
//
// Profiles: Represents instragram's user profile.
// Account:  Represents instagram's personal account.
// Search:   Represents instagram's search.
// Timeline: Represents instagram's timeline.
// Activity: Represents instagram's user activity.
// Inbox:    Represents instagram's messages.
// Location: Represents instagram's locations.
//
// See Scheme section in README.md for more information.
//
// We recommend to use Export and Import functions after first Login.
//
// Also you can use SetProxy and UnsetProxy to set and unset proxy.
// Golang also provides the option to set a proxy using HTTP_PROXY env var.
type Instagram struct {
	user string
	pass string
	// user-agent: Instagram 107.0.0.27.121 Android (24/7.0; 380dpi; 1080x1920; OnePlus; ONEPLUS A3010; OnePlus3T; qcom; en_US)
	userAgent string
	// device id: android-1923fjnma8123
	dID string
	// uuid: 8493-1233-4312312-5123
	uuid string
	// rankToken
	rankToken string
	// token
	token string
	// phone id
	pid string
	// SessionId
	//  This is a temporary ID which changes in the official app every time the
	//  user closes and re-opens the Instagram application or switches account.
	sid string
	// pigeonSessionId: 1bf37a34-86ff-5909-865c-0ff95082e698
	psid string
	// mid (from cookie)
	mid string
	// ads id
	adid string
	// challenge URL
	challengeURL string
	// auth token we receive
	auth string
	// wwwClaim we receive
	wwwClaim string
	// pwKeyId
	pwKeyId string
	// pwPubKey
	pwPubKey string

	// Instagram objects

	// Challenge controls security side of account (Like sms verify / It was me)
	Challenge *Challenge
	// Profiles is the user interaction
	Profiles *Profiles
	// Account stores all personal data of the user and his/her options.
	Account *Account
	// Search performs searching of multiple things (users, locations...)
	Search *Search
	// Timeline allows to receive timeline media.
	Timeline *Timeline
	// Activity are instagram notifications.
	Activity *Activity
	// Inbox are instagram message/chat system.
	Inbox *Inbox
	// Feed for search over feeds
	Feed *Feed
	// User contacts from mobile address book
	Contacts *Contacts
	// Location instance
	Locations *LocationInstance

	c *http.Client
}

// SetHTTPClient sets http client.  This further allows users to use this functionality
// for HTTP testing using a mocking HTTP client Transport, which avoids direct calls to
// the Instagram, instead of returning mocked responses.
func (inst *Instagram) SetHTTPClient(client *http.Client) {
	inst.c = client
}

// SetHTTPTransport sets http transport. This further allows users to tweak the underlying
// low level transport for adding additional fucntionalities.
func (inst *Instagram) SetHTTPTransport(transport http.RoundTripper) {
	inst.c.Transport = transport
}

// SetDeviceID sets device id
func (inst *Instagram) SetDeviceID(id string) {
	inst.dID = id
}

// SetUUID sets uuid
func (inst *Instagram) SetUUID(uuid string) {
	inst.uuid = uuid
}

// SetPhoneID sets phone id
func (inst *Instagram) SetPhoneID(id string) {
	inst.pid = id
}

// SetCookieJar sets the Cookie Jar. This further allows to use a custom implementation
// of a cookie jar which may be backed by a different data store such as redis.
func (inst *Instagram) SetCookieJar(jar http.CookieJar) error {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return err
	}
	// First grab the cookies from the existing jar and we'll put it in the new jar.
	cookies := inst.c.Jar.Cookies(url)
	inst.c.Jar = jar
	inst.c.Jar.SetCookies(url, cookies)
	return nil
}

// New creates Instagram structure
func New(username, password string) *Instagram {
	rand.Seed(time.Now().Unix())
	// this call never returns error
	jar, _ := cookiejar.New(nil)
	inst := &Instagram{
		user: username,
		pass: password,
		userAgent: generateUserAgent(),
		wwwClaim: "0",
		dID: generateDeviceID(
			generateMD5Hash(username + password),
		),
		uuid: generateUUID(), // both uuid must be differents
		pid:  generateUUID(),
		sid: generateUUID(),
		psid: generateUUID(),
		c: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Jar: jar,
		},
	}
	inst.init()

	return inst
}

func (inst *Instagram) init() {
	inst.Challenge = newChallenge(inst)
	inst.Profiles = newProfiles(inst)
	inst.Activity = newActivity(inst)
	inst.Timeline = newTimeline(inst)
	inst.Search = newSearch(inst)
	inst.Inbox = newInbox(inst)
	inst.Feed = newFeed(inst)
	inst.Contacts = newContacts(inst)
	inst.Locations = newLocation(inst)
}

// SetProxy sets proxy for connection.
func (inst *Instagram) SetProxy(url string, insecure bool) error {
	uri, err := neturl.Parse(url)
	if err == nil {
		inst.c.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		}
	}
	return err
}

// UnsetProxy unsets proxy for connection.
func (inst *Instagram) UnsetProxy() {
	inst.c.Transport = nil
}

// Save exports config to ~/.goinsta
func (inst *Instagram) Save() error {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("home") // for plan9
	}
	return inst.Export(filepath.Join(home, ".goinsta"))
}

// Export exports *Instagram object options
func (inst *Instagram) Export(path string) error {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return err
	}

	config := ConfigFile{
		ID:        inst.Account.ID,
		User:      inst.user,
		UserAgent: inst.userAgent,
		DeviceID:  inst.dID,
		UUID:      inst.uuid,
		RankToken: inst.rankToken,
		Token:     inst.token,
		Auth:      inst.auth,
		WWWClaim:  inst.wwwClaim,
		PhoneID:   inst.pid,
		Cookies:   inst.c.Jar.Cookies(url),
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)
}

// Export exports selected *Instagram object options to an io.Writer
func Export(inst *Instagram, writer io.Writer) error {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return err
	}

	config := ConfigFile{
		ID:        inst.Account.ID,
		User:      inst.user,
		UserAgent: inst.userAgent,
		DeviceID:  inst.dID,
		UUID:      inst.uuid,
		RankToken: inst.rankToken,
		Token:     inst.token,
		Auth:      inst.auth,
		WWWClaim:  inst.wwwClaim,
		PhoneID:   inst.pid,
		Cookies:   inst.c.Jar.Cookies(url),
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}

// ImportReader imports instagram configuration from io.Reader
//
// This function does not set proxy automatically. Use SetProxy after this call.
func ImportReader(r io.Reader) (*Instagram, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	config := ConfigFile{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return ImportConfig(config)
}

// ImportConfig imports instagram configuration from a configuration object.
//
// This function does not set proxy automatically. Use SetProxy after this call.
func ImportConfig(config ConfigFile) (*Instagram, error) {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return nil, err
	}

	inst := &Instagram{
		user:      config.User,
		userAgent: config.UserAgent,
		dID:       config.DeviceID,
		uuid:      config.UUID,
		rankToken: config.RankToken,
		token:     config.Token,
		pid:       config.PhoneID,
		auth:      config.Auth,
		wwwClaim:  config.WWWClaim,
		sid: 	   generateUUID(),
		psid: 	   generateUUID(),
		c: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
	inst.c.Jar, err = cookiejar.New(nil)
	if err != nil {
		return inst, err
	}
	inst.c.Jar.SetCookies(url, config.Cookies)

	for _, cookie := range inst.c.Jar.Cookies(url) {
		if strings.ToLower(cookie.Name) == "mid" {
			inst.mid = cookie.Value
		}
	}

	if inst.userAgent == "" {
		return inst, fmt.Errorf("user_agent missing in session file - please renew file")
	}
	if inst.auth == "" {
		return inst, fmt.Errorf("auth missing in session file - please renew file")
	}
	if inst.wwwClaim == "" {
		return inst, fmt.Errorf("www_claim missing in session file - please renew file")
	}

	inst.init()
	inst.Account = &Account{inst: inst, ID: config.ID}
	inst.Account.Sync()

	return inst, nil
}

// Import imports instagram configuration
//
// This function does not set proxy automatically. Use SetProxy after this call.
func Import(path string) (*Instagram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ImportReader(f)
}

func (inst *Instagram) readMsisdnHeader() error {
	data, err := json.Marshal(
		map[string]string{
			"device_id": inst.uuid,
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint:   urlMsisdnHeader,
			IsPost:     true,
			Connection: "keep-alive",
			Query:      generateSignature(b2s(data)),
		},
	)
	return err
}

func (inst *Instagram) contactPrefill() error {
	data, err := json.Marshal(
		map[string]string{
			"phone_id":   inst.pid,
			"_csrftoken": inst.token,
			"usage":      "prefill",
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint:   urlContactPrefill,
			IsPost:     true,
			Connection: "keep-alive",
			Query:      generateSignature(b2s(data)),
		},
	)
	return err
}

func (inst *Instagram) zrToken() error {
	_, err := inst.sendRequest(
		&reqOptions{
			Endpoint:   urlZrToken,
			IsPost:     false,
			Connection: "keep-alive",
			Query: map[string]string{
				"device_id":        inst.dID,
				"token_hash":       "",
				"custom_device_id": inst.uuid,
				"fetch_reason":     "token_expired",
			},
		},
	)
	return err
}

func (inst *Instagram) sendAdID() error {
	data, err := inst.prepareData(
		map[string]interface{}{
			"adid": inst.adid,
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint:   urlLogAttribution,
			IsPost:     true,
			Connection: "keep-alive",
			Query:      generateSignature(data),
		},
	)
	return err
}

// Login performs instagram login.
//
// Password will be deleted after login
func (inst *Instagram) Login() error {
	err := inst.readMsisdnHeader()
	if err != nil {
		return err
	}

	//err = inst.syncFeatures(goInstaExperiments)
	err = inst.syncFeatures(goInstaLoginExperiments)
	if err != nil {
		return err
	}

	err = inst.zrToken()
	if err != nil {
		return err
	}

	err = inst.sendAdID()
	if err != nil {
		return err
	}

	err = inst.contactPrefill()
	if err != nil {
		return err
	}

	enc_pass, err := GenerateEncPassword(inst.pass,inst.pwKeyId,inst.pwPubKey)
	if err != nil {
		return err
	}

	result, err := json.Marshal(
		map[string]interface{}{
			"guid":                inst.uuid,
			"login_attempt_count": 0,
			//"_csrftoken":          inst.token,
			"device_id":           inst.dID,
			"adid":                inst.adid,
			"phone_id":            inst.pid,
			"username":            inst.user,
			//"password":            inst.pass,
			"enc_password":		   enc_pass,
			"google_tokens":       "[]",
			"country_codes":       `{ country_code: "1", source: "default" }`,
			"jazoest":				GenerateJazoest(inst.pid),
		},
	)
	if err != nil {
		return err
	}
	body, err := inst.sendRequest(
		&reqOptions{
			Endpoint: urlLogin,
			Query:    generateSignature(b2s(result)),
			IsPost:   true,
			Login:    true,
			Connection: "close",
		},
	)
	if err != nil {
		return err
	}
	inst.pass = ""

	// getting account data
	res := accountResp{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	inst.Account = &res.Account
	inst.Account.inst = inst
	inst.rankToken = strconv.FormatInt(inst.Account.ID, 10) + "_" + inst.uuid
	inst.zrToken()

	return err
}

// Logout closes current session
func (inst *Instagram) Logout() error {
	_, err := inst.sendSimpleRequest(urlLogout)
	inst.c.Jar = nil
	inst.c = nil
	return err
}

func (inst *Instagram) syncFeatures(features string) error {
	data, err := inst.prepareData(
		map[string]interface{}{
			"id":          inst.uuid,
			"experiments": features,
		},
	)
	if err != nil {
		return err
	}

	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint: urlQeSync,
			Query:    generateSignature(data),
			IsPost:   true,
			Login:    true,
		},
	)
	return err
}

func (inst *Instagram) megaphoneLog() error {
	data, err := inst.prepareData(
		map[string]interface{}{
			"id":        inst.Account.ID,
			"type":      "feed_aysf",
			"action":    "seen",
			"reason":    "",
			"device_id": inst.dID,
			"uuid":      generateMD5Hash(string(time.Now().Unix())),
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint: urlMegaphoneLog,
			Query:    generateSignature(data),
			IsPost:   true,
			Login:    true,
		},
	)
	return err
}

func (inst *Instagram) expose() error {
	data, err := inst.prepareData(
		map[string]interface{}{
			"id":         inst.Account.ID,
			"experiment": "ig_android_profile_contextual_feed",
		},
	)
	if err != nil {
		return err
	}

	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint: urlExpose,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)

	return err
}

// GetMedia returns media specified by id.
//
// The argument can be int64 or string
//
// See example: examples/media/like.go
func (inst *Instagram) GetMedia(o interface{}) (*FeedMedia, error) {
	media := &FeedMedia{
		inst:   inst,
		NextID: o,
	}
	return media, media.Sync()
}
