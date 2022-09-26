package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const apikey = "PocyicaipdytNomeffyevUtyoaflebOs"
var homeport = 3000

func rando(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func rdpfinder(remoteip string, remoteport string) bool {
	magic := "\x03\x00\x00\x2c\x27\xe0\x00\x00\x00\x00\x00Cookie: mstshash=eltons\r\n\x01\x00\x08\x00\x00\x00\x00\x00"
	con, err := net.Dial("tcp", remoteip+":"+remoteport)
	if err != nil {
		return false
	} else {
		defer con.Close()
		fmt.Fprint(con, magic)
		reply := make([]byte, 19)
		_, err = con.Read(reply)
		response := hex.EncodeToString(reply)
		if len(response) == 38 && strings.Contains(response, "03000013") {
			return true
		}
	}
	return false
}

func forward(conn net.Conn, remoteip string, remoteport string) {
	client, err := net.Dial("tcp", remoteip+":"+remoteport)
	if err != nil {
		log.Println("Dial failed: %v", err)
	}
	log.Printf("Connected to localhost %v\n", conn)
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func listenup(homeip string, homeport int, remoteip string, remoteport string) {
	listener, err := net.Listen("tcp", homeip+":"+strconv.Itoa(homeport))
	if err != nil {
		log.Printf("Failed to setup listener after listenup goroutine called: %v", err)
		return
	}

	log.Printf("Listening on " + homeip + ":" + strconv.Itoa(homeport))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection %v\n", conn)
		go forward(conn, remoteip, remoteport)
	}
}

func backend() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/forward", forwardapi).Methods("POST")
	http.Handle("/", router)

	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatalf("SSL Server Error: " + err.Error())
		os.Exit(0)
	}
}

func forwardapi(response http.ResponseWriter, request *http.Request) {
	ip := strings.Split(request.RemoteAddr, ":")[0]
	log.Println(ip + ": /forward")
	request.ParseForm()

	apikeychallenge := request.FormValue("apikey")
	homeip := request.FormValue("homeip")
	remoteip := request.FormValue("remoteip")
	remoteport := request.FormValue("remoteport")

	if apikeychallenge == apikey {
		_, err := net.DialTimeout("tcp", remoteip+":"+remoteport, 8*time.Second)
		if err != nil {
			fmt.Fprintf(response, "CANT_CONNECT")
			return
		}
		if rdpfinder(remoteip, remoteport) == false {
			fmt.Fprintf(response, "NOT_RDP")
			return
		}
		listener, err := net.Listen("tcp", homeip+":"+strconv.Itoa(homeport))
		if err != nil {
			fmt.Fprintf(response, "LISTEN_ERROR")
		} else {
			listener.Close()
			go listenup(homeip, homeport, remoteip, remoteport)
			fmt.Fprintf(response, homeip+":"+strconv.Itoa(homeport))
			homeport++
		}
	} else {
		fmt.Fprintf(response, "BADAPIKEY")
	}
	return
}

func indexHandler(response http.ResponseWriter, request *http.Request) {
	ip := strings.Split(request.RemoteAddr, ":")[0]
	log.Println(ip + ": /")
	fmt.Fprintf(response, "")
}

func main() {
	fmt.Println("Server init...")
	go backend()
	for true {
		time.Sleep(100 * time.Millisecond)
	}
}
