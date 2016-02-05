package authserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"log"
	"net/http"
	"net/url"
	//"reflect"
)

func HandleLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	if r.Method == "POST" &&
		r.Form.Get("login") == "test" &&
		r.Form.Get("password") == "test" {
		return true
	}

	w.Write([]byte("<html><body>"))

	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s\" method=\"POST\">",
		ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri))))

	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

	return false
}

//func DownloadAccessToken(url string, auth *osin.BasicAuth, output map[string]interface{}) error {
func DownloadAccessToken(url string, auth interface{}, output map[string]interface{}) error {
	log.Println("DownloadAccessToken func")
	// download access token
	preq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if auth != nil {
		switch v := auth.(type) {
		case *osin.BasicAuth:
			log.Println("basic auth")
			preq.SetBasicAuth(v.Username, v.Password)
		case *osin.BearerAuth:
			// Direct Setting
			preq.Header.Set("Authorization", "Bearer "+v.Code)
			log.Println("bearer : ", v.Code)
		default:
			//var r = reflect.TypeOf(auth)
			log.Println("skip if another type")
		}
	}

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		return err
	}

	if presp.StatusCode != 200 {
		return errors.New("Invalid status code")
	}

	jdec := json.NewDecoder(presp.Body)
	err = jdec.Decode(&output)
	return err
}
