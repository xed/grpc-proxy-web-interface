package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func getPath() string {
	path, exists := os.LookupEnv("DUMP_PATH")
	if !exists {
		log.Fatal("no DUMP_PATH")
	}
	return path
}
func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("no PORT")
	}
	return port
}
func getLogin() string {
	if getAuth() == "true" {
		login, exists := os.LookupEnv("LOGIN")
		if !exists {
			log.Fatal("no LOGIN")
		}
		return login
	}
	return ""
}
func getPass() string {
	if getAuth() == "true" {
		pass, exists := os.LookupEnv("PASSWORD")
		if !exists {
			log.Fatal("no PASSWORD")
		}
		return pass
	}
	return ""
}
func getAuth() string {
	auth, exists := os.LookupEnv("AUTH")
	if !exists {
		log.Fatal("no AUTH")
	}
	return auth
}
func main() {
	path := getPath()
	port := getPort()
	user := getLogin()
	password := getPass()
	auth := getAuth()
	log.Print("Starting DUMP_PATH:", path, " PORT:", port, " AUTH:", auth, " USER:", user, " PASS:", password)
	if getAuth() == "true" {
		http.HandleFunc("/", GetOnly(BasicAuth(HandleIndex)))
	} else {
		http.HandleFunc("/", GetOnly(HandleIndex))
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func readDump() []PrintDump {
	path := getPath()
	log.Print("Starting DUMP:", getPath())
	var dump PrintDumps
	i := 0
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++
		text := scanner.Text()
		jsonDump, err := prettyPrint([]byte(text))
		if err != nil {
			log.Print("Error :", err)
		}
		dump = append(dump, PrintDump{
			Number: i,
			Json:   string(jsonDump),
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	sort.Sort(sort.Reverse(dump))
	return dump
}

type PrintDump struct {
	Number int
	Json   string
}
type PrintDumps []PrintDump

func (a PrintDumps) Len() int           { return len(a) }
func (a PrintDumps) Less(i, j int) bool { return a[i].Number < a[j].Number }
func (a PrintDumps) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func generateHTML() string {

	style := "<style>" +
		"* {box-sizing: border-box}" +
		".tab {float: left; border: 1px solid #ccc; background-color: #f1f1f1;  width: 20%;}" +
		".tab button {display: block; background-color: inherit; color: black;  padding: 22px 16px;width: 100%;border: none;outline: none;text-align: left;cursor: pointer;transition: 0.3s;font-size: 17px;}" +
		".tab button:hover {background-color: #ddd;}" +
		".tab button.active {background-color: #ccc;}" +
		".tabcontent {float: left; padding: 0px 12px; border: 1px solid #ccc;  width: 80%;  border-left: none;}" +
		"</style>"
	jsonString := "<!DOCTYPE html>" +
		"<html>" +
		"<head>" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">" +
		style +
		"</head>" +
		"<body>"
	dump := readDump()

	jsonString += "<div class=\"tab\">"
	for i, _ := range dump {
		id := "Dump" + strconv.Itoa(i)
		jsonString += "<button id=button" + strconv.Itoa(i) + " class=\"tablinks\" onclick=\"openDiv(event, '" + id + "')\">" + id + "</button>"
	}
	jsonString += "</div>"

	for i, printDump := range dump {

		id := "Dump" + strconv.Itoa(i)
		jsonString += "<div id=\"" + id + "\" class=\"tabcontent\">"
		jsonString += "<pre class=\"json-renderer\">" + printDump.Json + "</pre>"
		tabEnd := "</div>"
		jsonString += tabEnd
	}

	jsString := "<script>" +
		"function openDiv(evt, divName) {var i, tabcontent, tablinks; tabcontent = document.getElementsByClassName(\"tabcontent\"); for (i = 0; i < tabcontent.length; i++) { tabcontent[i].style.display = \"none\";}tablinks = document.getElementsByClassName(\"tablinks\"); for (i = 0; i < tablinks.length; i++) { tablinks[i].className = tablinks[i].className.replace(\" active\", \"\"); } document.getElementById(divName).style.display = \"block\";evt.currentTarget.className += \" active\";}" +
		"document.getElementById(\"button0\").click();" +
		"</script>"
	jsonString += jsString + "</body></html>"
	return jsonString
}
