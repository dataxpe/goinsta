package goinsta

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type reqOptions struct {
	// Connection is connection header. Default is "close".
	Connection string

	// Login process
	Login bool

	// Endpoint is the request path of instagram api
	Endpoint string

	// IsPost set to true will send request with POST method.
	//
	// By default this option is false.
	IsPost bool

	// UseV2 is set when API endpoint uses v2 url.
	UseV2 bool

	// Query is the parameters of the request
	//
	// This parameters are independents of the request method (POST|GET)
	Query map[string]string
}

func (insta *Instagram) sendSimpleRequest(uri string, a ...interface{}) (body []byte, err error) {
	return insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(uri, a...),
		},
	)
}

func (insta *Instagram) sendRequest(o *reqOptions) (body []byte, err error) {
	method := "GET"
	if o.IsPost {
		method = "POST"
	}
	if o.Connection == "" {
		o.Connection = "keep-alive"
	}

	nu := goInstaAPIUrl
	if o.UseV2 {
		nu = goInstaAPIUrlv2
	}

	u, err := url.Parse(nu + o.Endpoint)
	if err != nil {
		return nil, err
	}

	vs := url.Values{}
	bf := bytes.NewBuffer([]byte{})

	for k, v := range o.Query {
		vs.Add(k, v)
	}

	if o.IsPost {
		bf.WriteString(vs.Encode())
	} else {
		for k, v := range u.Query() {
			vs.Add(k, strings.Join(v, " "))
		}

		u.RawQuery = vs.Encode()
	}

	var req *http.Request
	req, err = http.NewRequest(method, u.String(), bf)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", insta.userAgent)
	req.Header.Set("X-Ads-Opt-Out", "0")
	req.Header.Set("X-DEVICE-ID", insta.uuid)
	req.Header.Set("X-CM-Bandwidth-KBPS", "-1.000")
	req.Header.Set("X-CM-Latency", "-1.000")
	req.Header.Set("X-IG-App-Locale", goInstaLanguage)
	req.Header.Set("X-IG-Device-Locale", goInstaLanguage)
	req.Header.Set("X-Pigeon-Session-Id", insta.psid)
	req.Header.Set("X-Pigeon-Rawclienttime", fmt.Sprintf("%d.000",time.Now().Unix())) // 1560144153.925
	req.Header.Set("X-IG-Connection-Speed", fmt.Sprintf("%dkbps", acquireRand(1000, 3700)))
	req.Header.Set("X-IG-Bandwidth-Speed-KBPS", "-1.000")
	req.Header.Set("X-IG-Bandwidth-TotalBytes-B", "0")
	req.Header.Set("X-IG-Bandwidth-TotalTime-MS", "0")
	req.Header.Set("X-IG-EU-DC-ENABLED", "0")
	req.Header.Set("X-IG-Extended-CDN-Thumbnail-Cache-Busting-Value", "1000")
	req.Header.Set("X-Bloks-Version-Id", goIInstaBloksVersion)
	req.Header.Set("X-MID", insta.mid)
	req.Header.Set("X-IG-WWW-Claim", insta.wwwClaim)
	req.Header.Set("X-Bloks-Is-Layout-RTL", "false")
	req.Header.Set("X-IG-Connection-Type", connType)
	req.Header.Set("X-IG-Capabilities", igCapabilities)
	req.Header.Set("X-IG-App-ID", fbAnalytics)
	req.Header.Set("X-IG-Device-ID", insta.uuid)
	req.Header.Set("X-IG-Android-ID", insta.dID)
	req.Header.Set("Accept-Language", strings.Replace(goInstaLanguage,"_","-",1))
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("Authorization", insta.auth)
	req.Header.Set("Connection", o.Connection)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	u, _ = url.Parse(goInstaAPIUrl)
	for _, value := range insta.c.Jar.Cookies(u) {
		if strings.Contains(value.Name, "csrftoken") {
			insta.token = value.Value
		} else if strings.ToLower(value.Name) == "mid" {
			insta.mid = value.Value
		}
	}

	if os.Getenv("GOINSTA_DEBUG") != "" {
		/*
		j, _ := json.MarshalIndent(req.Header,"< ","\t")
		fmt.Printf("< ---- SEND HEADER ----\n< %s\n< --------\n\n",j)
		*/


		dumprequest,_ := httputil.DumpRequest(req,false)
		fmt.Printf("< ----- SEND ----\n%s\n< -----------\n",dumprequest)
		j, _ := json.MarshalIndent(vs,"> ","\t")
		fmt.Printf("< ==== POST BODY ====\n< %s\n< =========\n\n\n\n",j)

		//j, _ = json.MarshalIndent(resp.Header,"> ","\t")
		//fmt.Printf("> ---- RECV HEADER ----\n> %s\n> --------\n\n",j)

		dumpresp ,_ := httputil.DumpResponse(resp,false)
		fmt.Printf("> ----- RECV ----\n%s\n",dumpresp)
	}

	if resp.Header.Get("ig-set-authorization") != "" {
		auth := resp.Header.Get("ig-set-authorization")
		if auth[len(auth)-1:] != ":" {
			insta.auth = auth
		}
	}
	if resp.Header.Get("x-ig-set-www-claim") != "" {
		insta.wwwClaim = resp.Header.Get("x-ig-set-www-claim")
	}
	if resp.Header.Get("ig-set-password-encryption-key-id") != "" {
		insta.pwKeyId = resp.Header.Get("ig-set-password-encryption-key-id")
	}
	if resp.Header.Get("ig-set-password-encryption-pub-key") != "" {
		insta.pwPubKey = resp.Header.Get("ig-set-password-encryption-pub-key")
	}

	// handle gzip response
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err = ioutil.ReadAll(reader)
	if err == nil {
		if os.Getenv("GOINSTA_DEBUG") != "" {
			fmt.Printf("> ==== RECV BODY ====\n> Status: %d %s\n>\n> %s\n> ======\n\n",
				resp.StatusCode, resp.Status, body)
		}

		err = isError(resp.StatusCode, body)
	}

	return body, err
}

func isError(code int, body []byte) (err error) {
	switch code {
	case 200:
	case 503:
		return Error503{
			Message: "Instagram API error. Try it later.",
		}
	case 400:
		ierr := Error400{}
		err = json.Unmarshal(body, &ierr)
		if err != nil {
			return err
		}

		if ierr.Message == "challenge_required" {
			return ierr.ChallengeError

		}

		if err == nil && ierr.Message != "" {
			return ierr
		}
	default:
		ierr := ErrorN{}
		err = json.Unmarshal(body, &ierr)
		if err != nil {
			return err
		}
		return ierr
	}
	return nil
}

func (insta *Instagram) prepareData(other ...map[string]interface{}) (string, error) {
	data := map[string]interface{}{
		"_uuid":      insta.uuid,
		"_csrftoken": insta.token,
	}
	if insta.Account != nil && insta.Account.ID != 0 {
		data["_uid"] = strconv.FormatInt(insta.Account.ID, 10)
	}

	for i := range other {
		for key, value := range other[i] {
			data[key] = value
		}
	}
	b, err := json.Marshal(data)
	if err == nil {
		return b2s(b), err
	}
	return "", err
}

func (insta *Instagram) prepareDataQuery(other ...map[string]interface{}) map[string]string {
	data := map[string]string{
		"_uuid":      insta.uuid,
		"_csrftoken": insta.token,
	}
	for i := range other {
		for key, value := range other[i] {
			data[key] = toString(value)
		}
	}
	return data
}

func acquireRand(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
