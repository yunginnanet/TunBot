package main

import (
	"fmt"
	"net/http"
	"net/url"
	"crypto/tls"
	
	"strings"
	"strconv"
	
	"bytes"
	
	"io/ioutil"
	
	"time"
	
	"encoding/hex"
	"encoding/base64"
	"crypto/md5"
	"github.com/sparrc/go-ping"
	
	"gopkg.in/AlecAivazis/survey.v1"
	
)

var (
	"127.0.0.1",
	}
)

const (
	apikey = "hpjGsW388tuUJruT"
	sslInsecureSkipVerify bool = true
	
	masterpass = "42f658b75e96614d035f66982bde4ad3"
	salt = "C4BbprF4X"
)

func httpclient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: sslInsecureSkipVerify},
	}
	client := &http.Client{Transport: tr}
	return client
}

func bestlatency() (string) {
	smallest := int64(99999999999999999)
	winner := "null"
	for i := 0; i < len(serverlist); i++ {
		target := serverlist[i]
		pinger, err := ping.NewPinger(target)
		if err != nil {
			continue
		}
		pinger.SetPrivileged(true)
		pinger.Count = 2
		pinger.Run() // blocks until finished
		stats := pinger.Statistics() // get send/receive/rtt stats
		ms := int64(stats.AvgRtt / time.Millisecond)
		fmt.Println(serverlist[i] + " - " + strconv.FormatInt(ms,10) + "ms")
		if ms < smallest {
			smallest = ms
			winner = serverlist[i]
		}
	}
	return winner
}

func apicall(tunserv string,remoteip string, remoteport string) {
	client := httpclient()
	
	homeip := "0.0.0.0"
	targeturl := "https://"+tunserv+"/forward"
	
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
	
	resp_body, _ := ioutil.ReadAll(resp.Body)
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
			if strings.Contains(response,":") {
				forwarded := strings.Split(response,":")
				fmt.Println("\nSuccess! You can now connect to Remote Desktop at "+tunserv+" on port "+forwarded[1])
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
    hasher.Write([]byte(passwordattempt+salt))
    passwordattempt = hex.EncodeToString(hasher.Sum(nil))
	if passwordattempt == masterpass {
		return true
	} else {
		fmt.Println("Authentication failure!")
		return false
	}
}


func banner() {
	clear, _ := base64.StdEncoding.DecodeString("ICAgIF9fICBfXyAgICAgICAgICAgICAgICAgIF9cbiAgIC8gLyAvIC9fICBfX19fX18gIF9fX18gIChfKSAgX19cbiAgLyAvXy8gLyAvIC8gLyBfXyBcLyBfXyBcLyAvIHwvXy9cbiAvIF9fICAvIC9fLyAvIC9fLyAvIC8gLyAvIC8+ICA8XG4vXy8gL18vXF9fLCAvXF9fX18vXy8gL18vXy9fL3xffFxuICAgICAgL19fX18v")
	bannerstring := string(clear)
	bannersplit := strings.Split(bannerstring,"\\n")
	for i := 0; i < len(bannersplit); i++ {
		fmt.Println(bannersplit[i])
	}
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
			Remoteip	string
			Remoteport	string
		}{}
		
		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		
		fmt.Print("\nFinding tunnel server with lowest latency...\n\n")
		tunserv := bestlatency()

		//tunserv = "127.0.0.1" //bypass for testing
		//fmt.Println("\nServer bypass enabled, using "+tunserv)
		
		fmt.Println("\nRequesting a tunnel to "+answers.Remoteip+":"+answers.Remoteport+" from "+tunserv+"...\n")
		
		
		
		apicall(tunserv,answers.Remoteip,answers.Remoteport)
	
	}
	
}
