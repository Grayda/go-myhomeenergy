package myhomeenergy

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

var token string        // Token is like a session key. A successful login returns a token you can use for the duration of your session
var Email string        // The email you're using to log in. Not sure if we really need to store this, but some functions might rely on it
var passwordHash string // The hash of your password. The hash is calculated by doing this: sha256(emailAddress + ":" + md5(password))
var LoggedIn bool       // We're successfully logged in.

// Command is what we GET to the server, in order to retrieve information
type Command struct {
	cmd  string
	data string
	fmt  string
}

type Data struct {
	Token       string
	NumAccounts string
}

// LoginResponse this is what we get back from the server when we log in. Not every field is filled, so check before using that data
type LoginResponse struct {
	Gip struct {
		Error       string
		Rc          string
		Version     string
		Numaccounts string
		Token       string
	}
}

// Response is a general response from the server, after login.
type Response struct {
	Gwrcmds struct {
		Gwrcmd struct {
			Gcmd  string
			Gdata struct {
				Gip struct {
					Goal    string
					Power   string
					Rc      string
					Version string
					Meter   []struct {
						Current  string
						Datatype string
						Did      string
						Enabled  string
						Error    string
						Factor   string
						Master   string
						Name     string
						Nodetype string
						Port     string
						Prodtype string
						Status   string
					}
					Chart struct {
						Costperhalfhour   string
						Energyperhalfhour string
					}
				}
			}
		}
	}
}

// Login logs you in to the myHomeEnergy website, given an email address and password.
// Nothing is stored, but if you don't want to pass in your details, consider using LoginToken instead
func Login(emailaddress string, password string) (string, error) {
	passwordHash = generateHash(emailaddress, password)
	token, err := LoginHash(emailaddress, passwordHash)
	return token, err
}

// LoginHash logs you in to the myHomeEnergy website, given a pre-calculated password.
// This token can be calculated like so: sha256(emailaddress:md5(password))
// The colon is a string literal, so don't forget that!
func LoginHash(emailaddress string, hash string) (string, error) {
	fmt.Println("Logging in with", emailaddress, "and", hash)
	data := &Command{
		cmd:  "GWRLogin",
		data: "<gip><version>1</version><email>" + emailaddress + "</email><password>" + hash + "</password></gip>",
	}

	jsonData, err := sendLoginRequest(data)

	switch jsonData.Gip.Rc {
	case "200":
		fmt.Println("Login OK!")
		info, _ := getMeters(jsonData.Gip.Token)
		chart := getInfo(jsonData.Gip.Token, info)
		spew.Dump(chart)
		return jsonData.Gip.Token, err
	default:
		return "", err
	}

}

func getInfo(token string, resp Response) Response {
	// t := time.Now()
	// td := fmtdate.Format("YYYYMMDD", t.AddDate(0, 0, 1)) + "000000"
	// fd := fmtdate.Format("YYYYMMDD", t) + "000000"

	data := &Command{
		cmd:  "GWRBatch",
		data: "<gwrcmds><gwrcmd><gcmd>DeviceGetChart</gcmd><gdata><gip><version>1</version><token>" + token + "</token><did>" + resp.Gwrcmds.Gwrcmd.Gdata.Gip.Meter[0].Did + "</did><period>now</period><feed>energyperhalfhour,costperhalfhour</feed></gip></gdata></gwrcmd><gwrcmd><gcmd>UserGetChart</gcmd><gdata><gip><version>1</version><token>" + token + "</token><period>now</period><feed>energyperhalfhour,costperhalfhour,tempout</feed></gip></gdata></gwrcmd></gwrcmds>",
	}

	info, _ := sendRequest(data)

	return info

}

func getMeters(token string) (Response, error) {
	data := &Command{
		cmd:  "GWRBatch",
		data: "<gwrcmds><gwrcmd><gcmd>SPA_UserGetSmartMeterList</gcmd><gdata><gip><version>1</version><token>" + token + "</token></gip></gdata></gwrcmd></gwrcmds>",
	}

	resp, err := sendRequest(data)

	return resp, err

}

func generateHash(emailaddress string, password string) string {
	tmp := []byte(password)
	mHash := fmt.Sprintf("%x", md5.Sum(tmp))
	tmp = []byte(emailaddress + ":" + mHash)
	hash := sha256.Sum256(tmp)
	return fmt.Sprintf("%x", hash)
}

func sendRequest(data *Command) (Response, error) {
	var jsonData Response

	r, _ := http.Get("https://myhomeenergy.com.au/gwr/gop.php?cmd=" + data.cmd + "&fmt=json&data=" + data.data)
	response, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(response, &jsonData)

	return jsonData, err
}

func sendLoginRequest(data *Command) (LoginResponse, error) {
	var jsonData LoginResponse
	r, _ := http.Get("https://myhomeenergy.com.au/gwr/gop.php?cmd=" + data.cmd + "&fmt=json&data=" + data.data)
	response, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(response))
	err := json.Unmarshal(response, &jsonData)
	return jsonData, err
}
