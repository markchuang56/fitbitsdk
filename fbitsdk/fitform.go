package fbitsdk

import (
	"bytes"
	"encoding/json"
	//"crypto"
	//"io"
	"io/ioutil"
	"net/http"
	"net/url"
	//"strings"
	"fmt"
)

const (
	CmdKey        = "cmdKey"
	CmdActivities = "activities"
	CmdSleep      = "sleep"
	CmdBodyFat    = "bodyfat"
	CmdHeartRate  = "heartrate"

	srcUrl1    = "https://api.fitbit.com/1/user/"
	activities = "/activities/date/"
	heartRate  = "/activities/heart/date/"
	bodyFat    = "/body/log/fat/date/"
	srcUrl2    = "https://api.fitbit.com/1.2/user/"
	sleep      = "/sleep/date/"

	//introspect := "https://api.fitbit.com/1.1/oauth2/introspect"

	//AuthURL:  "https://www.fitbit.com/oauth2/authorize",
	//TokenURL: "https://api.fitbit.com/oauth2/token",

	//RedirectURL: "http://localhost:8020/cb",
	//Scopes:      []string{"activity", "heartrate", "location", "nutrition", "profile", "settings", "sleep", "social", "weight"},

	//ClientID:     "22D6FQ",
	//ClientSecret: "be9c1fb74ca0d6b8c93deb35ba305093",
)

func sdkUrlPrepare(form url.Values, uid string) (string, error) {

	mx := form

	for k, v := range mx {
		fmt.Println(k)
		switch k {
		case CmdKey:
			fmt.Println(v[0])
			switch v[0] {
			case CmdActivities:
				sdkActivities(mx, uid)
				break

			case CmdSleep:
				sdkSleep(mx, uid)
				break

			case CmdBodyFat:
				sdkBodyFat(mx, uid)
				break

			case CmdHeartRate:
				sdkHeartRate(mx, uid)
				break

			default:
				break
			}
			return v[0], nil
			break

		default:
			fmt.Println("== CMD FAIL==")
			fmt.Println(k)
			break
		}
	}

	return "", nil
}

func fitDecode(resp *http.Response, data interface{}) {

	//if err != nil {
	//	fmt.Println("GET ERROR, 完了 ...")
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	fmt.Println(resp)
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	fmt.Println(newStr)
	fmt.Println("== HA HA ==")

	var fitData []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(fitData)
	fmt.Println("== NEW DECORDER ==")
	fmt.Println(fitData)
	//
	xbody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("=== READ ALL ERROR ===")
		panic(err.Error())
	}
	fmt.Println(xbody)
	fmt.Println()

	var dataUProfile UserProfile
	var dataUContent userContent
	var dataUtatal UserTatal
	var dataUtatal_x UserTatal_tx
	json.Unmarshal(xbody, &dataUProfile)
	json.Unmarshal(xbody, &dataUContent)
	json.Unmarshal(xbody, &dataUtatal)
	fmt.Printf("Results: %v\n", dataUProfile)
	fmt.Printf("Results: %v\n", dataUContent)
	fmt.Printf("Results: %v\n", dataUtatal)

	json.NewDecoder(resp.Body).Decode(dataUtatal_x)

	fmt.Println(dataUtatal_x)

	decoder := json.NewDecoder(resp.Body)
	//var data Tracks
	err = decoder.Decode(&dataUtatal_x)
	if err != nil {
		fmt.Println("JSON ERROR ...")
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(string(xbody[v.Offset-40 : v.Offset]))
		}
	}

	var dailies []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(dailies)
	fmt.Println(dailies)
	fmt.Println("=============")
	//
	fmt.Println(resp.Header)
	fmt.Println(resp.Body)
	fmt.Println(resp.ContentLength)
	fmt.Println(resp.TransferEncoding)
	fmt.Println(resp.Request)

	fmt.Println("USER DO ANYTHING ...")
}

/*

	//var body io.Reader
	//body = strings.NewReader(form.Encode())
	//req, err := http.NewRequest(http.MethodGet, urlStr, body)
	//if err != nil {
	//	return nil, err
	//}
	//ay := "abc"
	//return &ay, nil

type keyValue struct{ key, value []byte }

type byKeyValue []keyValue

func (p byKeyValue) Len() int      { return len(p) }
func (p byKeyValue) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p byKeyValue) Less(i, j int) bool {
	sgn := bytes.Compare(p[i].key, p[j].key)
	if sgn == 0 {
		sgn = bytes.Compare(p[i].value, p[j].value)
	}
	return sgn < 0
}

func (p byKeyValue) appendValues(values url.Values) byKeyValue {
	for k, vs := range values {
		k := encode(k, true)
		for _, v := range vs {
			v := encode(v, true)
			p = append(p, keyValue{k, v})
		}
	}
	return p
}

*/

/*

// encode encodes string per section 3.6 of the RFC. If double is true, then
// the encoding is applied twice.
func encode(s string, double bool) []byte {
	// Compute size of result.
	m := 3
	if double {
		m = 5
	}
	n := 0
	for i := 0; i < len(s); i++ {
		if noEscape[s[i]] {
			n++
		} else {
			n += m
		}
	}

	p := make([]byte, n)

	// Encode it.
	j := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		if noEscape[b] {
			p[j] = b
			j++
		} else if double {
			p[j] = '%'
			p[j+1] = '2'
			p[j+2] = '5'
			p[j+3] = "0123456789ABCDEF"[b>>4]
			p[j+4] = "0123456789ABCDEF"[b&15]
			j += 5
		} else {
			p[j] = '%'
			p[j+1] = "0123456789ABCDEF"[b>>4]
			p[j+2] = "0123456789ABCDEF"[b&15]
			j += 3
		}
	}
	return p
}
*/
