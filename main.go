package main

import (
    "net"
    "net/http"
    "net/url"
    "fmt"
    "flag"
    "log"
    "os"
    "bufio"
    "time"
    "strconv"
)

var headers = []string{"X-Originating-IP","X-Forwarded-For","X-Remote-IP","X-Remote-Addr","X-Client-IP","X-Host","X-Forwarded-Host","Origin","Host"}
var inject = []string{"127.0.0.1","localhost","0.0.0.0","0","127.1","127.0.1","2130706433"}
var urlt string
var pfile string
var to int

func payloadInject() {
    timeout := time.Duration(to * 1000000)
    var tr = &http.Transport{
                MaxIdleConns:      30,
                IdleConnTimeout:   time.Second,
                DisableKeepAlives: true,
                DialContext: (&net.Dialer{
                        Timeout:   timeout,
                        KeepAlive: time.Second,
                }).DialContext,
        }
        client := &http.Client{
                Transport:     tr,
                Timeout:       timeout,
        }
    file, err := os.Open(pfile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        for _, header := range headers {
            req, err := http.NewRequest("GET",urlt, nil)
            req.Header.Set(header, scanner.Text())
            resp, err := client.Do(req)
            if err != nil {
                continue
            }
            fmt.Println("[*] "+"["+urlt+"]"+" "+"["+header+": "+scanner.Text()+"]"+" "+" [Code: "+strconv.Itoa(int(resp.StatusCode))+"] "+"[Size: "+ strconv.Itoa(int(resp.ContentLength))+"]")
            defer resp.Body.Close()
        }
    }
}

func headerInject() {
    timeout := time.Duration(to * 1000000)
    var tr = &http.Transport{
                MaxIdleConns:      30,
                IdleConnTimeout:   time.Second,
                DisableKeepAlives: true,
                DialContext: (&net.Dialer{
                        Timeout:   timeout,
                        KeepAlive: time.Second,
                }).DialContext,
        }
        client := &http.Client{
                Transport:     tr,
                Timeout:       timeout,
        }
    for _, header := range headers {
        for _, i := range inject {
            req, err := http.NewRequest("GET",urlt, nil)
            req.Header.Set(header, i)
            resp, err := client.Do(req)
            if err != nil {
                continue
            }
            fmt.Println("[*] "+"["+urlt+"]"+" "+"["+header+": "+i+"]"+" "+" [Code: "+strconv.Itoa(int(req.StatusCode))+"] "+"[Size: "+ strconv.Itoa(int(req.ContentLength))+"]")
            defer resp.Body.Close()
        }
    }
}

func main() {
    flag.StringVar(&urlt, "url", "", "target URL")
    flag.StringVar(&pfile, "pfile","","payload file")
    flag.IntVar(&to, "t", 10000, "timeout (milliseconds)")
    flag.Parse()
    if urlt == "" {
        flag.PrintDefaults()
    } else {
        u, err := url.Parse(urlt)
        if err != nil {
            log.Fatal(err)
        }
        if u.Scheme == "" || u.Host == "" || u.Path == "" {
            fmt.Println("Invalid URL: ",urlt)
            os.Exit(1)
        }
        if pfile != "" {
            payloadInject()
        } else {
            headerInject()
        }
    }
}
