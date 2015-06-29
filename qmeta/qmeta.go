package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fileloc := flag.String("file", "config", "config file contains RabbitMQ node addrs.")
	flag.Parse()

	currentTime := time.Now().Format("20060102150405")
	desDir := "qmeta_" + currentTime
	os.MkdirAll(desDir, 0744)

	f, err := os.Open(*fileloc)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout: timeout}

	c := make(chan string, 50)
	cnt := 0
	buffer := bufio.NewReader(f)
	for {
		line, err := buffer.ReadString('\n')
		line = strings.TrimSpace(line)
		if err == io.EOF {
			if len(line) <= 0 {
				break
			}
		} else if err != nil {
			fmt.Printf("%s\n", err.Error())
			break
		}

		cfg := strings.Split(line, " ")
		fieldLen := len(cfg)
		if fieldLen != 3 {
			fmt.Printf("config format error. Expect \"Host User Passwd\", but receive \"%s\".\n", line)
			continue
		}
		host := cfg[0]
		user := cfg[1]
		passwd := cfg[2]

		cnt++
		go download(client, desDir, host, user, passwd, c)
	}
	for i := 0; i < cnt; i++ {
		<-c
		//h := <-c
		//fmt.Printf("main receive: %s\n", h)
	}
	close(c)
}

func download(client *http.Client, desDir string, host, user, passwd string, c chan<- string) {
	defer func() {
		c <- host
	}()

	auth := []byte(user + ":" + passwd)
	base64edAuth := base64.StdEncoding.EncodeToString(auth)

	url := "http://" + host + "/api/definitions?download=rabbit.json&auth=" + base64edAuth
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("err ocur.", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err ocur.", err.Error())
		return
	}

	filename := strings.Split(host, ":")[0]
	metaf, err := os.OpenFile(desDir+"/"+filename+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer metaf.Close()
	fmt.Printf("Host %s backup done!\n", filename)
	metaf.WriteString(string(body[:]))
}
