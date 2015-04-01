package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type hit_msg struct {
	Hostname   string `json:"hostname"`
	Smooth_num string `json:"smooth_num"`
	Smoothness string `json:"smoothness"`
}

type Client struct {
	Hostname   string `json:"hostname"`
	Init       string `json:"init"`
	Inc        string `json:"inc"`
	Last_itr   string `json:"last_itr"`
	Smooth_num string `json:"smooth_num"`
	Smoothness string `json:"smoothness"`
}

func sendUpdate() {
	im := Client{
		Hostname:   *hostname,
		Init:       *start,
		Inc:        *inc,
		Last_itr:   last_itr,
		Smoothness: smoothest,
		Smooth_num: smooth_num,
	}
	jsonBytes, err := json.Marshal(&im)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(*server+"/itr", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
		return
	}
	resp.Body.Close()
}

func sendHit() {
	hm := hit_msg{
		Hostname:   *hostname,
		Smooth_num: smooth_num,
		Smoothness: smoothest,
	}
	jsonBytes, err := json.Marshal(&hm)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(*server+"/hit", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
		return
	}
	resp.Body.Close()
}

var hostname *string
var smoothest string
var smooth_num string
var start *string
var inc *string
var last_itr string
var server *string

func main() {
	start = flag.String("start", "113851762", "Start exponent")
	inc = flag.String("inc", "1000000000", "Increment amount")
	server = flag.String("server", "http://localhost:3000", "Server to report to")
	hostname = flag.String("name", "", "client name, default is hostname")
	flag.Parse()

	if *hostname == "" {
		*hostname, _ = os.Hostname()
	}

	// Create the command that will be run
	smoothFinder := exec.Command("./ssearch", *start, *inc)

	// Get the output pipe of the command
	out, err := smoothFinder.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(out)

	// Notify the server we are starting a job
	sendUpdate()

	// Start the command
	err = smoothFinder.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Read the output and send data to the server
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		// Remove trailing newline
		line = line[:len(line)-1]
		// Split on :
		cmdv := strings.Split(line, ":")
		if len(cmdv) < 1 {
			continue
		}

		switch cmdv[0] {
		case "ITR":
			last_itr = cmdv[1]
			log.Println(last_itr)
			sendUpdate()
			break
		case "HIT":
			smoothest = cmdv[1]
			smooth_num = cmdv[2]
			sendHit()
			break
		default:
			continue
		}
	}
}
