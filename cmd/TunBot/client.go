package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"gopkg.in/AlecAivazis/survey.v1"
)

var (
	apikey = os.Getenv("TUNBOT_API_KEY")

	maintenance                = os.Getenv("TUNBOT_MAINTENANCE")
	salt                       = os.Getenv("TUNBOT_SALT")
	sslInsecureSkipVerify bool = true
)

func init() {
	for _, val := range []*string{&apikey, &maintenance, &salt} {
		if *val == "" {
			*val = "fuckmeupfam"
		}
	}
}

func httpclient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: sslInsecureSkipVerify},
	}
	client := &http.Client{Transport: tr}
	return client
}

func apicall(tunserv string, remoteip string, remoteport string) {
	client := httpclient()

	homeip := "0.0.0.0"
	targeturl := "https://" + tunserv + "/forward"

	data := url.Values{}

	data.Set("apikey", apikey)
	data.Set("homeip", homeip)
	data.Set("remoteip", remoteip)
	data.Set("remoteport", remoteport)

	u, _ := url.ParseRequestURI(targeturl)

	urlStr := fmt.Sprintf("%v", u)
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println("Cannot connect to tunnel server!")
		return
	}

	resp_body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode == 200 {
		response := string(resp_body)

		switch response {

		case "CANT_CONNECT":
			fmt.Println("Remote server cannot connect to target!")
			break

		case "NOT_RDP":
			fmt.Println("No RDP server detected on target!")
			break

		case "LISTEN_ERROR":
			fmt.Println("Listen error from remote server...")
			break

		case "BADAPIKEY":
			fmt.Println("Bad API key detected! This will be reported.")
			break

		default:
			if strings.Contains(response, ":") {
				forwarded := strings.Split(response, ":")
				fmt.Println("\nSuccess! You can now connect to Remote Desktop at " + tunserv + " on port " + forwarded[1])
			} else {
				fmt.Println("Received malformed response!")
			}
		}
	}
}

func authenticate() bool {
	passwordattempt := ""
	prompt := &survey.Password{
		Message: "Password",
	}
	survey.AskOne(prompt, &passwordattempt, nil)

	hasher := md5.New()
	hasher.Write([]byte(passwordattempt + salt))
	passwordattempt = hex.EncodeToString(hasher.Sum(nil))
	if passwordattempt == maintenance {
		return true
	} else {
		fmt.Println("Authentication failure!")
		return false
	}
}

func banner() {
	clear, _ := base64.StdEncoding.DecodeString("ICAgIF9fICBfXyAgICAgICAgICAgICAgICAgIF9cbiAgIC8gLyAvIC9fICBfX19fX18gIF9fX18gIChfKSAgX19cbiAgLyAvXy8gLyAvIC8gLyBfXyBcLyBfXyBcLyAvIHwvXy9cbiAvIF9fICAvIC9fLyAvIC9fLyAvIC8gLyAvIC8+ICA8XG4vXy8gL18vXF9fLCAvXF9fX18vXy8gL18vXy9fL3xffFxuICAgICAgL19fX18v")
	io.Copy(os.Stdout, bytes.NewReader(clear))
	intro := [...]string{
		"-----------RDP Booster-----------",
	}
	for i := 0; i < len(intro); i++ {
		fmt.Println(intro[i])
	}
}

func main() {

	authenticated := authenticate()
	if authenticated == true {
		banner()

		var qs = []*survey.Question{
			{
				Name:     "remoteip",
				Prompt:   &survey.Input{Message: "What is the Remote Desktop Address?"},
				Validate: survey.Required,
			},
			{
				Name:     "remoteport",
				Prompt:   &survey.Input{Message: "What port is Remote Desktop running on?"},
				Validate: survey.Required,
			},
		}

		answers := struct {
			Remoteip   string
			Remoteport string
		}{}

		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Print("\nFinding tunnel server with lowest latency...\n\n")

		tunserv := "127.0.0.1" //bypass for testing

		fmt.Println("\nRequesting a tunnel to " + answers.Remoteip + ":" + answers.Remoteport + " from " + tunserv + "...\n")

		apicall(tunserv, answers.Remoteip, answers.Remoteport)

	}

}
