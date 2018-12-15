package main

import (
	//"bytes"
	//"context"
	"fmt"

	"encoding/base64"
	//"encoding/gob"
	//"encoding/json"

	//"github.com/gorilla/sessions"

	//"golang.org/x/oauth2"
	//"golang.org/x/oauth2/fitbit"
	//"io/ioutil"

	//"errors"
	//"io"
	//"strings"

	//"log"
	"net/http"
	//"net/url"
	//"time"
	fbitsdk "./fbitsdk"
	"html/template"
	"os"
)

var (
	homeTmpl          = template.Must(template.New("home").ParseFiles("templates/table.html"))
	homeLoggedOutTmpl = template.Must(template.New("loggout").ParseFiles("templates/loggedout.html"))
)

func init() {

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, `<html><body><a href="/auth/fitbit">Sign in with Fitbit</a></body></html>`)
	w.Header().Set("Content-Type", "text/html; charset-utf-8") //
	if err := homeLoggedOutTmpl.ExecuteTemplate(w, "loggedout.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func xxhandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><a href="/auth/fitbit">Sign in with Fitbit</a></body></html>`)
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	fbitsdk.HandleAuth(w, r)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== FIIBIT DEMO CALLBACK ===")
	token, err := fbitsdk.HandleAuthCallback(w, r)
	if err != nil {
		fmt.Println("== Callback Fail ==")
	}

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

	http.Redirect(w, r, "/user", http.StatusFound)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("== demo user comes ...")
	w.Header().Set("Content-Type", "text/html; charset-utf-8") //
	if err := homeTmpl.ExecuteTemplate(w, "table.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fmt.Println("=== CALLBACK ADDR ===")
	//fmt.Println(oauthConfig.RedirectURL)

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
	http.HandleFunc("/cb", handleAuthCallback)
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
}

func demoServeGetSleep(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET SLEEP ====")

	fmt.Println("==== SERVE GET SLEEP ==== OK ")
}

func demoServeGetHeartRate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET ACTIVITIES ====")

}

func demoServeGetActivities(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET ACTIVITIES ====")

	fmt.Println("==== SERVE GET ACTIVITIES ==== OK ")
}

func demoServeGetUserProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET USER PROFILE ====")

	fmt.Println("==== SERVE GET USER PROFILE ==== OK ")
}

func demoServeGetUserBodyWeight(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==== SERVE GET USER PROFILE ====")

	fmt.Println("==== SERVE GET USER BODY WEIGHT ==== OK ")
}
