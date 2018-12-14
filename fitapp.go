package main

import (
	"bytes"
	"context"
	"fmt"

	"encoding/base64"
	"encoding/gob"
	"encoding/json"

	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
	"io/ioutil"

	//"errors"
	"io"
	"strings"

	//"log"
	"net/http"
	"net/url"
	//"time"
	"html/template"
	"os"
)

var (
	store = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_SECRET")), []byte(os.Getenv("COOKIE_ENCRYPT")))

	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("FITBIT_OAUTH_ID"),
		ClientSecret: os.Getenv("FITBIT_OAUTH_SECRET"),
		Endpoint:     fitbit.Endpoint,

		Scopes:      []string{"activity", "heartrate", "location", "nutrition", "profile", "settings", "sleep", "social", "weight"},
		RedirectURL: "https://fathomless-shore-18884.herokuapp.com/callback",

		//Endpoint: oauth2.Endpoint{
		//	AuthURL:  "https://www.fitbit.com/oauth2/authorize",
		//	TokenURL: "https://api.fitbit.com/oauth2/token",
		//},
		// See https://devcenter.heroku.com/articles/oauth#scopes
		//RedirectURL: "https://" + os.Getenv("HEROKU_APP_NAME") + "herokuapp.com/auth/heroku/callback", // See
		//
		//ClientID:     "22DD2F",
		//ClientSecret: "a62ee79d8e9ab5b3f6e99c6a775a16b5",
		//RedirectURL:  "http://127.0.0.1:8020/callback",
	}

	stateToken = os.Getenv("FITBIT_APP_NAME")

	homeTmpl          = template.Must(template.New("home").ParseFiles("templates/sleep.html"))
	homeLoggedOutTmpl = template.Must(template.New("loggout").ParseFiles("templates/loggedout.html"))
)

func init() {
	gob.Register(&oauth2.Token{})

	store.MaxAge(60 * 60 * 8)
	store.Options.Secure = true
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, `<html><body><a href="/auth/fitbit">Sign in with Fitbit</a></body></html>`)
	w.Header().Set("Content-Type", "text/html; charset-utf-8") //
	if err := homeLoggedOutTmpl.ExecuteTemplate(w, "loggedout.html", nil); err != nil {
		//if err := homeTmpl.ExecuteTemplate(w, "sleep.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("state Token")
	//fmt.Println(store["COOKIE_SECRET"])
	//fmt.Println(store["COOKIE_ENCRYPT"])
	fmt.Println(stateToken)
	url := oauthConfig.AuthCodeURL(stateToken)
	fmt.Println("=== OAUTH ===")
	fmt.Println(url)
	fmt.Println("=== URL for Authorize ===")
	http.Redirect(w, r, url, http.StatusFound)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== FIIBIT CALLBACK ===")

	if v := r.FormValue("state"); v != stateToken {
		http.Error(w, "Invalid State token", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	fmt.Println(ctx)
	fmt.Println()
	//fmt.Println(r)
	fmt.Println()
	fmt.Println(r.FormValue("code"))
	token, err := oauthConfig.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "fitbit-oauth-go")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(session)
	session.Values["fitbit-oauth-token"] = token
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var xt oauth2.Token
	//json.Unmarshal(token, &xt)
	fmt.Println(session)
	//testStruct := MyStruct{"hello world"}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(token)
	err = json.Unmarshal(reqBodyBytes.Bytes(), &xt)
	//reqBodyBytes.Bytes() // this is the []byte
	//if err := json.Unmarshal(byt, &obj); err != nil {
	//    panic(err)
	//}
	if err != nil {
		fmt.Println("=== ERROR ===")
		panic(err)
	}

	fmt.Println("=== FIIBIT CALLBACK 2 ===")

	fmt.Println("=== TOKEN ===")
	fmt.Println(reqBodyBytes.Bytes())
	fmt.Println(xt)
	fmt.Println(reqBodyBytes.String())
	//fmt.Println(token["AccessToken"])
	fmt.Println()

	fmt.Println()
	fmt.Println("AccessToken")
	fmt.Println(token.AccessToken)
	fmt.Println()

	fmt.Println("TokenType")
	fmt.Println(token.TokenType)
	fmt.Println()

	fmt.Println("RefreshToken")
	fmt.Println(token.RefreshToken)
	fmt.Println()

	fmt.Println("Expiry")
	fmt.Println(token.Expiry)
	fmt.Println()

	fmt.Println("GOOD!!")
	//http.Redirect(w, r, "/user", http.StatusFound)
	http.Redirect(w, r, "/userProfile", http.StatusFound)
	//http.Redirect(w, r, "/sleep", http.StatusFound)
	//http.Redirect(w, r, "/heartRate", http.StatusFound)
	//http.Redirect(w, r, "/activities", http.StatusFound)
	//http.Redirect(w, r, "/bodyWeight", http.StatusFound)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "fitbit-oauth-go")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, ok := session.Values["fitbit-oauth-token"].(*oauth2.Token)
	if !ok {
		http.Error(w, "Unable to assert token", http.StatusInternalServerError)
		return
	}

	fmt.Println(token.AccessToken)
	fmt.Println()

	//
	ctx := context.Background()

	//
	//var buf bytes.Buffer
	//buf.WriteString(oauthConfig.Endpoint.TokenURL)
	v := url.Values{
		//"token": {token.AccessToken}, //The OAuth 2.0 token to retrieve the state
		"token": {token.AccessToken},
	}
	//v.Set("redirect_uri", "https://fathomless-shore-18884.herokuapp.com/user")

	var body io.Reader
	//if r.method != http.MethodGet {
	body = strings.NewReader(v.Encode())
	//}
	fmt.Println("=== SAME or NOT ===")
	fmt.Println(body)

	//uProfile := "https://api.fitbit.com/1/user/6Z29KN/profile.json"
	//req, err := http.NewRequest(http.MethodGet, uProfile, body)

	introspect := "https://api.fitbit.com/1.1/oauth2/introspect"
	req, err := http.NewRequest(http.MethodPost, introspect, body)

	if err != nil {
		//return nil, err
		fmt.Println("NEW REQUEST ERROR !!!")
	}

	fmt.Println("== REQ ==")
	fmt.Println(req)
	fmt.Println(" VALUE")
	fmt.Println(req.Form)
	fmt.Println(req.Body)
	/*

		//for k, v := range c.Header {
		//	req.Header[k] = v
		//}
		fmt.Println(req.Header)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		fmt.Println(req.Header)

		req.Header.Set("client_id", "22DD2F")

		//req.Header.Set("redirect_uri", "https://fathomless-shore-18884.herokuapp.com/callback")
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		req.Header.Set("Authorization", "Basic Y2xpZW50X2lkOmNsaWVudCBzZWNyZXQ=")
	*/
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	xtype := token.TokenType
	xtype += " "
	xtype += token.AccessToken
	req.Header.Set("Authorization", xtype)
	fmt.Println(xtype)
	//token=<The OAuth 2.0 token to retrieve the state>
	//req.Header.Add("token", "The OAuth 2.0 token to retrieve the state")
	//req.Header.Set("grant_type", "authorization_code")

	//req.Header.Set("code", r.FormValue("code"))

	req.WithContext(ctx)
	fmt.Println(req)
	fmt.Println(req.Header)
	//

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Do(req)
	fmt.Println("CLIENT ...", resp.Body)
	fmt.Println(client.Transport)
	//GET https://api.fitbit.com/1/user/[user-id]/profile.json
	//user-id	The encoded ID of the user. Use "-" (dash) for current logged-in user.

	//uProfile := "https://api.fitbit.com/1/user/6Z29KN/profile.json"
	//resp, err := client.Get(uProfile)
	if err != nil {
		fmt.Println("GET ERROR, 完了 ...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	//var data interface{}
	//json.Unmarshal(resp.Body, &data)
	//fmt.Println(data)
	//fmt.Println(resp.Body.Reader)

	w.Header().Set("Content-Type", "text/html; charset-utf-8") //

	if err := homeTmpl.ExecuteTemplate(w, "sleep.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("USER DO ANYTHING ...")

}

func main() {
	fmt.Println("=== CALLBACK ADDR ===")
	fmt.Println(oauthConfig.RedirectURL)

	// Base 64
	//str := "c29tZSBkYXRhIHdpdGggACBhbmQg77u/"
	str := "Y2xpZW50X2lkOmNsaWVudCBzZWNyZXQ="
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%q\n", data)
	fmt.Println("==== BASE64 ====")

	http.HandleFunc("/", handleRoot)

	http.HandleFunc("/auth/fitbit", handleAuth)
	http.HandleFunc("/callback", handleAuthCallback)
	http.HandleFunc("/user", handleUser)

	http.HandleFunc("/sleep", demoServeGetSleep)
	http.HandleFunc("/heartRate", demoServeGetHeartRate)
	http.HandleFunc("/activities", demoServeGetActivities)
	http.HandleFunc("/userProfile", demoServeGetUserProfile)
	http.HandleFunc("/bodyWeight", demoServeGetUserBodyWeight)

	fmt.Println("FITBIT ...")
	port := os.Getenv("PORT")
	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "8020"
	}
	fmt.Println(port)
	http.ListenAndServe(":"+port, nil)
	//http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func demoServeGetSleep(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET SLEEP ====")
	session, err := store.Get(r, "fitbit-oauth-go")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, ok := session.Values["fitbit-oauth-token"].(*oauth2.Token)
	if !ok {
		http.Error(w, "Unable to assert token", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	var body io.Reader
	v := url.Values{
		"token": {token.AccessToken},
	}
	body = strings.NewReader(v.Encode())
	fmt.Println("=== SAME or NOT ===")
	fmt.Println(body)

	xSleep := "https://api.fitbit.com/1.2/user/6Z29KN/sleep/date/2018-12-03.json"
	req, err := http.NewRequest(http.MethodGet, xSleep, body)

	if err != nil {
		//return nil, err
		fmt.Println("NEW REQUEST ERROR !!!")
	}

	//fmt.Println("== REQ ==")
	//fmt.Println(req)
	//fmt.Println(" VALUE")
	//fmt.Println(req.Form)
	//fmt.Println(req.Body)

	xtype := token.TokenType
	xtype += " "
	xtype += token.AccessToken
	req.Header.Set("Authorization", xtype)
	fmt.Println(xtype)

	req.WithContext(ctx)
	fmt.Println(req)
	fmt.Println(req.Header)
	//

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Do(req)
	fmt.Println("CLIENT ...", resp.Body)
	fmt.Println(client.Transport)

	//resp, err := client.Get(uProfile)
	if err != nil {
		fmt.Println("GET ERROR, 完了 ...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(resp)
	defer resp.Body.Close()

	fmt.Println("=== GOOD ===")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	fmt.Println(newStr)
	fmt.Println("== HA HA ==")

	fmt.Println("==== SERVE GET SLEEP ==== OK ")
}

func demoServeGetHeartRate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET ACTIVITIES ====")
	session, err := store.Get(r, "fitbit-oauth-go")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, ok := session.Values["fitbit-oauth-token"].(*oauth2.Token)
	if !ok {
		http.Error(w, "Unable to assert token", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	var body io.Reader
	v := url.Values{
		"token": {token.AccessToken},
	}
	body = strings.NewReader(v.Encode())
	fmt.Println("=== SAME or NOT ===")
	fmt.Println(body)

	xheartrate := "https://api.fitbit.com/1/user/6Z29KN/activities/heart/date/today/1d.json"
	req, err := http.NewRequest(http.MethodGet, xheartrate, body)

	//introspect := "https://api.fitbit.com/1.1/oauth2/introspect"
	//req, err := http.NewRequest(http.MethodPost, introspect, body)

	if err != nil {
		//return nil, err
		fmt.Println("NEW REQUEST ERROR !!!")
	}

	fmt.Println("== REQ ==")
	fmt.Println(req)
	fmt.Println(" VALUE")
	fmt.Println(req.Form)
	fmt.Println(req.Body)

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	xtype := token.TokenType
	xtype += " "
	xtype += token.AccessToken
	req.Header.Set("Authorization", xtype)
	fmt.Println(xtype)

	req.WithContext(ctx)
	fmt.Println(req)
	fmt.Println(req.Header)
	//

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Do(req)
	fmt.Println("CLIENT ...", resp.Body)
	fmt.Println(client.Transport)

	//resp, err := client.Get(uProfile)
	if err != nil {
		fmt.Println("GET ERROR, 完了 ...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(resp)
	defer resp.Body.Close()

	fmt.Println("=== GOOD ===")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	fmt.Println(newStr)
	fmt.Println("== HA HA ==")

	fmt.Println("==== SERVE GET ACTIVITIES ==== OK ")
	//fmt.Fprintln(w, "GET HEART RATE")
}

func demoServeGetActivities(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET ACTIVITIES ====")
	session, err := store.Get(r, "fitbit-oauth-go")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, ok := session.Values["fitbit-oauth-token"].(*oauth2.Token)
	if !ok {
		http.Error(w, "Unable to assert token", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	var body io.Reader
	v := url.Values{
		"token": {token.AccessToken},
	}
	body = strings.NewReader(v.Encode())
	fmt.Println("=== SAME or NOT ===")
	fmt.Println(body)

	xactivities := "https://api.fitbit.com/1/user/6Z29KN/activities/steps/date/today/1m.json"
	//uProfile := "https://api.fitbit.com/1/user/6Z29KN/profile.json"
	req, err := http.NewRequest(http.MethodGet, xactivities, body)

	//introspect := "https://api.fitbit.com/1.1/oauth2/introspect"
	//req, err := http.NewRequest(http.MethodPost, introspect, body)

	if err != nil {
		//return nil, err
		fmt.Println("NEW REQUEST ERROR !!!")
	}

	fmt.Println("== REQ ==")
	fmt.Println(req)
	fmt.Println(" VALUE")
	fmt.Println(req.Form)
	fmt.Println(req.Body)

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	xtype := token.TokenType
	xtype += " "
	xtype += token.AccessToken
	req.Header.Set("Authorization", xtype)
	fmt.Println(xtype)

	req.WithContext(ctx)
	fmt.Println(req)
	fmt.Println(req.Header)
	//

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Do(req)
	fmt.Println("CLIENT ...", resp.Body)
	fmt.Println(client.Transport)

	//resp, err := client.Get(uProfile)
	if err != nil {
		fmt.Println("GET ERROR, 完了 ...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(resp)
	defer resp.Body.Close()

	fmt.Println("=== GOOD ===")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	fmt.Println(newStr)
	fmt.Println("== HA HA ==")
	//minlite.ApiGetActivities(w, r, cred, " ", " ")
	fmt.Println("==== SERVE GET ACTIVITIES ==== OK ")
}

func demoServeGetUserProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET USER PROFILE ====")
	session, err := store.Get(r, "fitbit-oauth-go")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, ok := session.Values["fitbit-oauth-token"].(*oauth2.Token)
	if !ok {
		http.Error(w, "Unable to assert token", http.StatusInternalServerError)
		return
	}

	//fmt.Println(token.AccessToken)
	//fmt.Println()

	//
	ctx := context.Background()

	var body io.Reader
	v := url.Values{
		"token": {token.AccessToken},
	}
	body = strings.NewReader(v.Encode())
	fmt.Println("=== SAME or NOT ===")
	fmt.Println(body)

	//uProfile := "https://api.fitbit.com/1/user/6Z29KN/profile.json"
	uProfile := "https://api.fitbit.com/1/user/-/profile.json"
	req, err := http.NewRequest(http.MethodGet, uProfile, body)

	if err != nil {
		//return nil, err
		fmt.Println("NEW REQUEST ERROR !!!")
	}

	fmt.Println("== REQ ==")
	fmt.Println(req)
	fmt.Println(" VALUE")
	fmt.Println(req.Form)
	fmt.Println(req.Body)

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	xtype := token.TokenType
	xtype += " "
	xtype += token.AccessToken
	req.Header.Set("Authorization", xtype)
	fmt.Println(xtype)

	req.WithContext(ctx)
	fmt.Println(req)
	fmt.Println(req.Header)
	//

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Do(req)
	fmt.Println("CLIENT ...", resp.Body)
	fmt.Println(client.Transport)

	//resp, err := client.Get(uProfile)
	if err != nil {
		fmt.Println("GET ERROR, 完了 ...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(resp)
	defer resp.Body.Close()

	fmt.Println("=== GOOD ===")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	fmt.Println(newStr)
	fmt.Println("== HA HA ==")

	if err := homeTmpl.ExecuteTemplate(w, "sleep.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("==== SERVE GET USER PROFILE ==== OK ")
}

func demoServeGetUserBodyWeight(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET USER PROFILE ====")
	session, err := store.Get(r, "fitbit-oauth-go")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, ok := session.Values["fitbit-oauth-token"].(*oauth2.Token)
	if !ok {
		http.Error(w, "Unable to assert token", http.StatusInternalServerError)
		return
	}

	//
	ctx := context.Background()

	var body io.Reader
	v := url.Values{
		"token": {token.AccessToken},
	}
	body = strings.NewReader(v.Encode())
	fmt.Println("=== SAME or NOT ===")
	fmt.Println(body)

	//uProfile := "https://api.fitbit.com/1/user/6Z29KN/profile.json"
	xbodyWeight := "https://api.fitbit.com/1/user/6Z29KN/body/log/fat/date/2018-12-03.json"
	req, err := http.NewRequest(http.MethodGet, xbodyWeight, body)

	if err != nil {
		//return nil, err
		fmt.Println("NEW REQUEST ERROR !!!")
	}

	fmt.Println("== REQ ==")
	fmt.Println(req)
	fmt.Println(" VALUE")
	fmt.Println(req.Form)
	fmt.Println(req.Body)

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	xtype := token.TokenType
	xtype += " "
	xtype += token.AccessToken
	req.Header.Set("Authorization", xtype)
	fmt.Println(xtype)

	req.WithContext(ctx)
	fmt.Println(req)
	fmt.Println(req.Header)
	//

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Do(req)
	fmt.Println("CLIENT ...", resp.Body)
	fmt.Println(client.Transport)

	//resp, err := client.Get(uProfile)
	if err != nil {
		fmt.Println("GET ERROR, 完了 ...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(resp)
	defer resp.Body.Close()

	fmt.Println("=== GOOD ===")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	fmt.Println(newStr)
	fmt.Println("== HA HA ==")

	if err := homeTmpl.ExecuteTemplate(w, "sleep.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("==== SERVE GET USER BODY WEIGHT ==== OK ")
}

type UserTatal struct {
	user userContent
}

//
type UserProfile struct {
	user string
	/*
		"user": {
		        "aboutMe":<value>,
		        "avatar":<value>,
		        "avatar150":<value>,
		        "avatar640":<value>,
		        "city":<value>,
		        "clockTimeDisplayFormat":<12hour|24hour>,
		        "country":<value>,
		        "dateOfBirth":<value>,
		        "displayName":<value>,
		        "distanceUnit":<value>,
		        "encodedId":<value>,
		        "foodsLocale":<value>,
		        "fullName":<value>,
		        "gender":<FEMALE|MALE|NA>,
		        "glucoseUnit":<value>,
		        "height":<value>,
		        "heightUnit":<value>,
		        "locale":<value>,
		        "memberSince":<value>,
		        "offsetFromUTCMillis":<value>,
		        "startDayOfWeek":<value>,
		        "state":<value>,
		        "strideLengthRunning":<value>,
		        "strideLengthWalking":<value>,
		        "timezone":<value>,
		        "waterUnit":<value>,
		        "weight":<value>,
		        "weightUnit":<value>
		    }
	*/
}

type UserTatal_tx struct {
	user UserProfile_tx
}

type UserProfile_tx struct {
	user string `json:"user"`
}

type userContent_tx struct {
	aboutMe                string `json:"aboutMe"`                //":<value>,
	avatar                 string `json:"avatar"`                 //:<value>,
	avatar150              string `json:"avatar150"`              ////":<value>,
	avatar640              string `json:"avatar640"`              //":<value>,
	city                   string `json:"city"`                   //":<value>,
	clockTimeDisplayFormat string `json:"clockTimeDisplayFormat"` //":<12hour|24hour>,
	country                string `json:"country"`                //":<value>,
	dateOfBirth            string `json:"dateOfBirth"`            //":<value>,
	displayName            string `json:"displayName"`            //":<value>,
	distanceUnit           string `json:"distanceUnit"`           //":<value>,
	encodedId              string `json:"encodedId"`              //":<value>,
	foodsLocale            string `json:"foodsLocale"`            //":<value>,
	fullName               string `json:"fullName"`               //":<value>,
	gender                 string `json:"gender"`                 //":<FEMALE|MALE|NA>,
	glucoseUnit            string `json:"glucoseUnit"`            //":<value>,
	height                 string `json:"height"`                 //":<value>,
	heightUnit             string `json:"heightUnit"`             //":<value>,
	locale                 string `json:"locale"`                 //":<value>,
	memberSince            string `json:"memberSince"`            //":<value>,
	offsetFromUTCMillis    string `json:"offsetFromUTCMillis"`    //":<value>,
	startDayOfWeek         string `json:"startDayOfWeek"`         //":<value>,
	state                  string `json:"state"`                  //":<value>,
	strideLengthRunning    string `json:"strideLengthRunning"`    //":<value>,
	strideLengthWalking    string `json:"strideLengthWalking"`    //":<value>,
	timezone               string `json:"timezone"`               //":<value>,
	waterUnit              string `json:"waterUnit"`              //":<value>,
	weight                 string `json:"weight"`                 //":<value>,
	weightUnit             string `json:"weightUnit"`             //":<value>
}

type userContent struct {
	aboutMe                string //":<value>,
	avatar                 string //:<value>,
	avatar150              string ////":<value>,
	avatar640              string //":<value>,
	city                   string //":<value>,
	clockTimeDisplayFormat string //":<12hour|24hour>,
	country                string //":<value>,
	dateOfBirth            string //":<value>,
	displayName            string //":<value>,
	distanceUnit           string //":<value>,
	encodedId              string //":<value>,
	foodsLocale            string //":<value>,
	fullName               string //":<value>,
	gender                 string //":<FEMALE|MALE|NA>,
	glucoseUnit            string //":<value>,
	height                 string //":<value>,
	heightUnit             string //":<value>,
	locale                 string //":<value>,
	memberSince            string //":<value>,
	offsetFromUTCMillis    string //":<value>,
	startDayOfWeek         string //":<value>,
	state                  string //":<value>,
	strideLengthRunning    string //":<value>,
	strideLengthWalking    string //":<value>,
	timezone               string //":<value>,
	waterUnit              string //":<value>,
	weight                 string //":<value>,
	weightUnit             string //":<value>
}

// Old
//func ApiGetEpochs(w http.ResponseWriter, r *http.Request, cred *ApiCredentials, tmStart, tmEnd string) []map[string]interface{} {
//	if loopBusy {
//		return nil
//	}
//	loopBusy = true
//	apiTimeStamp()
//	credx := oauth.Credentials{cred.Token, cred.Secret}
//	return serveGetEpochs(w, r, &credx)
//}

//GET https://api.fitbit.com/1/user/[user-id]/activities/date/[date].json
//user-id	The encoded ID of the user. Use "-" (dash) for current logged-in user.
//date	The date in the format yyyy-MM-dd

//GET https://api.fitbit.com/1/user/[user-id]/body/log/fat/date/[date].json
//GET https://api.fitbit.com/1/user/[user-id]/body/log/fat/date/[date]/[period].json
//GET https://api.fitbit.com/1/user/[user-id]/body/log/fat/date/[base-date]/[end-date].json
//user-id	The encoded ID of the user. Use "-" (dash) for current logged-in user.
//date	The date in the format yyyy-MM-dd.
//base-date	The end date when period is provided; range start date when a date range is provided. In the format yyyy-MM-dd or today.
//period	The date range period. One 1d, 7d, 1w, 1m.
//end-date	Range end date when date range is provided. Note: The range should not be longer than 31 days.

//GET https://api.fitbit.com/1/user/[user-id]/profile.json
//user-id	The encoded ID of the user. Use "-" (dash) for current logged-in user.

func serveGetEpochs(w http.ResponseWriter, r *http.Request) []map[string]interface{} {
	fmt.Println("==== SERVE GET EPOCHS ====")

	var sleeps []map[string]interface{}
	fmt.Println(sleeps)

	return sleeps
}

/* */

var HTTPClient contextKey

type contextKey struct{}

func contextClient(ctx context.Context) *http.Client {
	if ctx != nil {
		if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok && hc != nil {
			return hc
		}
	}
	return http.DefaultClient
}

type Client struct {
	Config oauth2.Config
	Header http.Header
}

type Credentials struct {
}

type request struct {
	method        string
	u             *url.URL
	form          url.Values
	verifier      string
	sessionHandle string
	callbackURL   string
}

/*
func (c *Client) do(ctx context.Context, urlStr string, r *request) (*http.Response, error) {
	//var body io.Reader

		if r.method != http.MethodGet {
			body = strings.NewReader(r.form.Encode())
		}
		req, err := http.NewRequest(r.method, urlStr, body)
		if err != nil {
			return nil, err
		}
		if req.URL.RawQuery != "" {
			return nil, errors.New("oauth: url must not contain a query string")
		}
		for k, v := range c.Header {
			req.Header[k] = v
		}
		r.u = req.URL
		auth, err := c.authorizationHeader(r)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", auth)
		if r.method == http.MethodGet {
			req.URL.RawQuery = r.form.Encode()
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		req = requestWithContext(ctx, req)
		fmt.Println(r)
		fmt.Println(req)
		fmt.Println("=== LIB ===")
		fmt.Println(req)

	client := contextClient(ctx)
	return client.Do(req)
}

// Get issues a GET to the specified URL with form added as a query string.
func (c *Client) Get(client *http.Client, credentials *Credentials, urlStr string, form url.Values) (*http.Response, error) {
	ctx := context.WithValue(context.Background(), HTTPClient, client)
	return c.GetContext(ctx, credentials, urlStr, form)
}

// GetContext uses Context to perform Get.
func (c *Client) GetContext(ctx context.Context, credentials *Credentials, urlStr string, form url.Values) (*http.Response, error) {
	return c.do(ctx, urlStr, &request{method: http.MethodGet, credentials: credentials, form: form})
}

// Post issues a POST with the specified form.
func (c *Client) Post(client *http.Client, credentials *Credentials, urlStr string, form url.Values) (*http.Response, error) {
	ctx := context.WithValue(context.Background(), HTTPClient, client)
	return c.PostContext(ctx, credentials, urlStr, form)
}

// PostContext uses Context to perform Post.
func (c *Client) PostContext(ctx context.Context, credentials *Credentials, urlStr string, form url.Values) (*http.Response, error) {
	return c.do(ctx, urlStr, &request{method: http.MethodPost, credentials: credentials, form: form})
}

*/

/*
// apiGet issues a GET request to the Twitter API and decodes the response JSON to data.
func apiGet(cred *oauth.Credentials, urlStr string, form url.Values, data interface{}) error {
	fmt.Println(" *****  API GET  *****")
	fmt.Println(urlStr)
	fmt.Println()
	resp, err := spxClient.Get(nil, cred, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func apiPost(cred *oauth.Credentials, urlStr string, form url.Values, data interface{}) error {
	fmt.Println(" *****  API POST  *****")
	resp, err := spxClient.Post(nil, cred, urlStr, form)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	fmt.Println(" *****  DECODE RESPONSE  *****")
	if resp.StatusCode != 200 {
		p, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("get %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
	}
	// TEST ...
	fmt.Println(resp.Body)
	return json.NewDecoder(resp.Body).Decode(data)
}
*/

/*

var homeTmpl = template.Must(template.New("home").ParseFiles("templates/epoch.html"))
var homeLoggedOutTmpl = template.Must(template.New("loggout").ParseFiles("templates/loggedout.html"))

func fitbitUserServeAuthorize(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" === FITBIT USER SERVE AUTHORIZE ===")
	ctx := context.Background()
	//
		conf := &oauth2.Config{
			ClientID:     "22D6FQ",
			ClientSecret: "be9c1fb74ca0d6b8c93deb35ba305093",
			Scopes:       []string{"SLEEP"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.fitbit.com/oauth2/authorize",
				TokenURL: "https://api.fitbit.com/oauth2/token",
			},
		}
	//

	conf := &oauth2.Config{
		ClientID:     "22DD2F",
		ClientSecret: "a62ee79d8e9ab5b3f6e99c6a775a16b5",
		//Scopes:       []string{"SLEEP", "activity"},
		Scopes: []string{"activity"},
		//Endpoint: oauth2.Endpoint{
		//	AuthURL:  "https://www.fitbit.com/oauth2/authorize",
		//	TokenURL: "https://api.fitbit.com/oauth2/token",
		//},

		//ClientID:     os.Getenv("HEROKU_OAUTH_ID"),
		//ClientSecret: os.Getenv("HEROKU_OAUTH_SECRET"),
		Endpoint: fitbit.Endpoint,
	}
	//
				22DD2F
				a62ee79d8e9ab5b3f6e99c6a775a16b5
				http://127.0.0.1:8080
				https://www.fitbit.com/oauth2/authorize
				https://api.fitbit.com/oauth2/token

		https://www.fitbit.com/oauth2/authorize?response_type=token&client_id=22DD2F&redirect_uri=http%3A%2F%2F127.0.0.1%3A8080&scope=activity%20heartrate%20location%20nutrition%20profile%20settings%20sleep%20social%20weight&expires_in=604800

	//

	xcallback := oauth2.SetAuthURLParam("redirect_uri", "http://127.0.0.1:8080/callback")
	//xcallback := oauth2.SetAuthURLParam("redirect_uri", "https://app-settings.fitbitdevelopercontent.com/simple-redirect.html")
	xtimeout := oauth2.SetAuthURLParam("expires_in", "325800")
	//xresponse := oauth2.SetAuthURLParam("response_type", "token")
	//fmt.Println()
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	//url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	url := conf.AuthCodeURL("state", xcallback, xtimeout)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	http.Redirect(w, r, url, http.StatusFound)

	fmt.Println()
	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	fmt.Println("STEP 1")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	fmt.Println("STEP 2")
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("STEP 3")
	client := conf.Client(ctx, tok)
	//client.Get("...")
	fmt.Println(client)

	//redir := oauthClient.AuthorizationURL(tempCred, nil)
	fmt.Println("=== TOKEN 1 ===")
	//fmt.Println(redir)
	fmt.Println("=== TOKEN 2 ===")

	//http.Redirect(w, r, redir, 302)

	//fitbit-app-224008
}

// authHandler reads the auth cookie and invokes a handler with the result.
type authHandler struct {
	//handler  func(w http.ResponseWriter, r *http.Request, c *oauth.Credentials)
	handler  func(w http.ResponseWriter, r *http.Request)
	optional bool
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" *****  SERVE HTTP  *****")
	//cred, _ := session.Get(r)[tokenCredKey].(*oauth.Credentials)
	//if cred == nil && !h.optional {
	//	http.Error(w, "Not logged in.", 403)
	//}

	h.handler(w, r)
}
*/

/*
// response responds to a request by executing the html remplate t with data.
func respond(w http.ResponseWriter, t *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println(" *****  RESPOND  *****")
	if err := t.Execute(w, data); err != nil {
		log.Print(err)
	}
}
*/

/*
func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" *****  SERVE HOME  *****")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	//fmt.Println(cred)
	fmt.Println("HOME")

	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	if err := homeLoggedOutTmpl.ExecuteTemplate(w, "loggedout.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//if cred == nil {
	//respond(w, homeLoggedOutTmpl, nil)
	//respond(w, homePage, nil)
	// for TEST
	//respond(w, homeTmpl, nil)
	//} else {
	//respond(w, homeTmpl, nil)
	//}
}

func fitbitConfig() {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "22D6FQ",
		ClientSecret: "be9c1fb74ca0d6b8c93deb35ba305093",
		Scopes:       []string{"SCOPE1", "SCOPE2"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.fitbit.com/oauth2/authorize",
			TokenURL: "https://api.fitbit.com/oauth2/token",
		},
	}
	fmt.Println(ctx)
	fmt.Println(conf)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!, HOW DO YOU DO ??\n")
}
*/

/*
func main() {
	//fitbitConfig()
	//ctx := context.Background()
	//
		data, err := ioutil.ReadFile("settings/index.json")
		fmt.Println("what ??")
		fmt.Println(data)
		fmt.Println("YES")
		//fmt.Println(ctx)
		if err != nil {
			fmt.Println(err)
		}
	//

	//
	//helloHandler := func(w http.ResponseWriter, req *http.Request) {
	//	io.WriteString(w, "Hello, world! what happen ??\n")
	//}

	http.Handle("/", &authHandler{handler: serveHome, optional: true})

	http.HandleFunc("/authorize", fitbitUserServeAuthorize)
	http.HandleFunc("/hello", helloHandler)

	errX := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if errX != nil {
		panic(errX)
	}

	//http.ListenAndServe(":8080", nil)

	//http.HandleFunc("/authorize", serveAuthorize)
	//log.Fatal(http.ListenAndServe(":8010", nil))
	//
		s := &http.Server{
			Addr:           ":8080",
			Handler:        helloHandler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Fatal(s.ListenAndServe())
	//
}



*/

//session, err := store.Get(r, "fitbit-oauth-go")
//fmt.Println(session)
//if err != nil {
//	fmt.Println("SESSION ERROR ...")
//	http.Error(w, err.Error(), http.StatusInternalServerError)
//	return
//}

// save the token
//session.Set(“AccessToken”, token.AccessToken)
//session.Set(“RefreshToken”, token.RefreshToken)
//session.Set(“TokenType”, token.TokenType)
//session.Set(“Expiry”, token.Expiry.Format(time.RFC3339))
//session.Save()
//client := conf.Client(oauth2.NoContext, tok)session.Values["fitbit-oauth-token"] = token

//if err := session.Save(r, w); err != nil {
//	fmt.Println("SESSION VALUE ERROR ...")
//	http.Error(w, err.Error(), http.StatusInternalServerError)
//	return
//}

/*
	url := oauthConfig.AuthCodeURL(stateToken)
	client := oauthConfig.Client(ctx, token)

	//tokenURL := "https://api.fitbit.com/oauth2/token"
	contentType := "application/x-www-form-urlencoded"
	xmethod := strings.NewReader("name=cjb")
	resp, err := client.Post(url, contentType, xmethod)

	fmt.Println(resp)
	if err != nil {
		fmt.Println("POST ERR ...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GOOD ... TOKEN")
	defer resp.Body.Close()
*/
//http.Redirect(w, r, "/user", http.StatusFound)

/*
	fmt.Println(token.TokenType)
	fmt.Println()

	fmt.Println(token.RefreshToken)
	fmt.Println()


	//
	urlStr := "https://api.fitbit.com/oauth2/token"
	var body io.Reader
	r.method = http.MethodPost
	//if r.method != http.MethodGet {
		body = strings.NewReader(r.form.Encode())
		//}
	req, err := http.NewRequest(r.method, urlStr, body)
	if err != nil {
		return nil, err
	}
	//if req.URL.RawQuery != "" {
	//	return nil, errors.New("oauth: url must not contain a query string")
	//}
	for k, v := range c.Header {
		req.Header[k] = v
	}
	//r.u = req.URL
	//auth, err := c.authorizationHeader(r)
	//if err != nil {
	//	return nil, err
	//}
	//req.Header.Set("Authorization", auth)
	//if r.method == http.MethodGet {
	//	req.URL.RawQuery = r.form.Encode()
	//} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//}
	req = requestWithContext(ctx, req)
	fmt.Println(r)
	fmt.Println(req)
	fmt.Println("=== LIB ===")
	fmt.Println(req)
	client := contextClient(ctx)
	return client.Do(req)

	client := oauthConfig.Client(context.Background(), token)
	fmt.Println("USER")
	fmt.Println(client)
	fmt.Println("what ???")

	//func (c *Client) Post(url, contentType string, body io.Reader) (resp *Response, err error)
	tokenURL := "https://api.fitbit.com/oauth2/token"
	contentType := "application/x-www-form-urlencoded"
	xmethod := strings.NewReader("name=cjb")
	resp, err := client.Post(tokenURL, contentType, xmethod)

	//resp, err := client.Get("https://api.heroku.com/account")
	fmt.Println(resp)
	if err != nil {
		fmt.Println("POST ERR ...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GOOD ...")
	defer resp.Body.Close()
*/
/*
	d := json.NewDecoder(resp.Body)
	var account struct { // See https://devcenter.heroku.com/articles/platform-api-reference#account
		Email string `json:"email"`
	}
	if err := d.Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `<html><body><h1>Hello Fitbit Demo %s</h1></body></html>`, account.Email)
*/
