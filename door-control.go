
package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
        "log"
	"net/http"
	"time"
)

var (
	// Use mcu pin 17, corresponds to physical pin 11 on the pi
	pinOutImpuls = rpio.Pin(17)
	pinInDoorOpen = rpio.Pin(26)
)

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pins
	//Pin feeding transistor - when High, transistor is on and impuls is sent - button is pushed
	pinOutImpuls.Output()
	pinOutImpuls.Low()

	//Pin reading door opened sensor
        // Low - door closed
        // High - door opened
        pinInDoorOpen.PullUp()
	pinInDoorOpen.Input()

	http.HandleFunc("/", indexFn);
	http.HandleFunc("/door", doorFn)
	http.HandleFunc("/pushButton", pushButtonFn)

	log.Printf("Starting http listener");
	log.Fatal(http.ListenAndServe(":8123", nil))
}

func doorFn(respWriter http.ResponseWriter, req *http.Request){
        log.Printf("doorFn %s",req.Method);
	if req.Method == "OPTIONS" {
		cors(respWriter)
		return
	}
	//fmt.Fprintf(respWriter,"door\n")
	state:=pinInDoorOpen.Read()
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	if state == rpio.High {
		log.Printf("door state: opened")
		fmt.Fprintf(respWriter,`{ "doorState" : "opened" }`);
		
	} else {
		log.Printf("door state: closed")
		fmt.Fprintf(respWriter,`{ "doorState" : "closed" }`);
	}
}

func pushButtonFn(respWriter http.ResponseWriter, req *http.Request){
        log.Printf("pushButtonFn %s", req.Method);
	if req.Method == "OPTIONS" {
		cors(respWriter)
		return
	}
	//fmt.Fprintf(respWriter,"pushButton\n")
	pinOutImpuls.High()
	log.Printf("button pushed")
	tmr := time.NewTimer(time.Second / 2 )
	go func(){
		<-tmr.C
		pinOutImpuls.Low()
		log.Printf("Button released")
	}()	
}

func indexFn(respWriter http.ResponseWriter, req *http.Request){
        log.Printf("indexFn %s", req.Method);
	//if req.Method == "OPTIONS" {
	//	cors(respWriter)
	//	return
	//}
	respWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	//fmt.Fprintf(respWriter,"<h1>Hi fropm Pi Zero W</h1>")
	fmt.Fprintf(respWriter,indexContent())
}

func cors(respWriter http.ResponseWriter){
	log.Printf("serving OPTIONS");
	respWriter.Header().Set("Access-Control-Allow-Origin", "*") 
	respWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	//respWriter.Header().Set("Access-Control-Max-Age: 86400 
	respWriter.WriteHeader(http.StatusOK)
}


func indexContent() string{
return `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Zero W</title>
</head>
<body>
<h1>Pi Zero W</h1>
Stav vrat: <span id="doorState">NA</span>
<br><br>
<button id="btnUpdateDoor">Update door state!</button>
<button id="btnPushButton">Push door button!</button>
<script>
document.addEventListener('DOMContentLoaded', onLoad, false);

function onLoad(ev){
	console.log("Page loaded");
	document.querySelector("#btnUpdateDoor").addEventListener("click",updateDoor);
	document.querySelector("#btnPushButton").addEventListener("click",pushButton);
	setInterval(()=>{updateDoor();}, 5000);
	updateDoor();
}
let lastDoorState="NA";
function updateDoor(){
	console.log("updating door state");
	var req=new XMLHttpRequest()
	req.open("GET","http://zero:8123/door")
	req.addEventListener("load",(resp)=>{
	  console.log("loaded:",req.responseText);
	  let respObj=JSON.parse(req.responseText);
	  document.querySelector("#doorState").textContent=respObj.doorState;
          if (lastDoorState != respObj.doorState){
            lastDoorState=respObj.doorState;
            window.navigator.vibrate(500)
	  }
	})

	req.addEventListener("error",(resp)=>{
	  console.log("error:",req.status)
	})

	req.send();
}

function pushButton(){
	console.log("pushing button");
	var req=new XMLHttpRequest()
	req.open("POST","http://zero:8123/pushButton")
	req.addEventListener("load",(resp)=>{
	  console.log("loaded:",req.responseText);
	})

	req.addEventListener("error",(resp)=>{
	  console.log("error:",req.status)
	})

	req.send();
}


</script>
</body>
</html>
`
}
