package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func getTID() (string, error) {
	var r http.Request
	r.ParseForm()
	fp := map[string]string{"os": "2", "browser": "Chrome86,0,4240,183", "fonts": "undefined", "screenInfo": "1440*900*24", "plugins": "Portable Document Format::internal-pdf-viewer::Chrome PDF Plugin|::mhjfbmdgcfjbbpaeojofohoefgiehjai::Chrome PDF Viewer|::internal-nacl-plugin::Native Client"}
	r.Form.Add("cb", "gen_callback")
	r.Form.Add("fp", createKeyValuePairs(fp))
	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest("POST", "https://passport.weibo.com/visitor/genvisitor", strings.NewReader(bodystr))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	etags := resp.Header["Etag"]
	if len(etags) == 0 {
		return "", fmt.Errorf("Etag empty")
	}
	return etags[0], nil
}

func GetSub(tid string) (string, error) {
	var cookieString string
	cookies := []string{}

	req, err := http.NewRequest("GET", "https://passport.weibo.com/visitor/visitor", nil)
	if err != nil {
		return cookieString, err
	}

	var resp *http.Response
	q := req.URL.Query()
	q.Add("a", "incarnate")
	q.Add("w", "2")
	q.Add("c", "095")
	q.Add("cb", "cross_domain")
	q.Add("from", "weibo")
	q.Add("t", url.QueryEscape(tid))
	q.Add("_rand", fmt.Sprintf("0.%08v%08v", rand.Int31n(100000000), rand.Int31n(100000000)))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return cookieString, err
	}
	for _, cookie := range resp.Cookies() {
		cookies = append(cookies, cookie.String())
	}
	cookieString = strings.Join(cookies, "; ")
	defer resp.Body.Close()
	return cookieString, nil
}

func main() {
	tid, err := getTID()
	if err == nil {
		fmt.Println(tid)
	}
	ss, _ := GetSub(tid)
	header := http.Header{}
	header.Add("Cookie", ss)
	request := http.Request{Header: header}

	fmt.Println(request.Cookies())
}
