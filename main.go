package myhomeenergy

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Token is the hash of our email address and password used to log in to myHomeEnergy
var Token string
var Email string
var LoggedIn bool

type Command struct {
	cmd  string
	data string
	fmt  string
}

// Login logs you in to the myHomeEnergy website, given an email address and password.
// If you don't want to pass in your details, consider using LoginToken instead
func Login(emailaddress string, password string) (bool, string, error) {
	Token = generateToken(emailaddress, password)
	LoginToken(emailaddress, Token)
	return true, Token, nil
}

// LoginToken logs you in to the myHomeEnergy website, given a pre-calculated token.
// This token can be calculated like so: sha256(emailaddress:md5(password))
// The colon is a literal one, so it'd look like this: sha256(abc@xyz.com:ABCDEF1234567890)
func LoginToken(emailaddress string, token string) (bool, error) {
	fmt.Println("Logging in with", emailaddress, "and", token)
	data := &Command{
		cmd:  "GWRLogin",
		data: "<gip><version>1</version><email>" + emailaddress + "</email><password>" + token + "</password></gip>",
		fmt:  "json",
	}
	// buf, _ := xml.Marshal(data)
	//
	// body := bytes.NewBuffer(buf)
	r, _ := http.Get("https://myhomeenergy.com.au/gwr/gop.php?cmd=" + data.cmd + "&fmt=" + data.fmt + "&data=" + data.data)
	response, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(response))
	return true, nil
}

func getMeters() {

}

func generateToken(emailaddress string, password string) string {
	tmp := []byte(password)
	mHash := fmt.Sprintf("%x", md5.Sum(tmp))
	tmp = []byte(emailaddress + ":" + mHash)
	hash := sha256.Sum256(tmp)
	return fmt.Sprintf("%x", hash)
}
