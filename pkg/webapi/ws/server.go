package main

import (
	"flag"
	"fmt"
	"github.com/weibaohui/mesh/pkg/webapi/ui"
 	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	go func() {
		for {
			if _, _, err := c.NextReader(); err != nil {
				c.Close()
				break
			}
		}
	}()
	logs := make(chan []byte)

	go ui.LogReader(logs)

	ticker := time.NewTicker(time.Second * 5)
	for {

		select {
		case item, ok := <-logs:
			fmt.Println("ok",ok)
			if !ok {
				break
			}
			nextWriter, err := c.NextWriter(websocket.TextMessage)
			defer nextWriter.Close()
			_, err = nextWriter.Write(item)
			if err != nil {
				log.Println("write:", err)
				break
			}
		case <-ticker.C:
			fmt.Println("ticker")
			nextWriter, err := c.NextWriter(websocket.TextMessage)
			defer nextWriter.Close()
			_, err = nextWriter.Write([]byte("ticker"))
			if err != nil {
				log.Println("write:", err)
				break
			}
			// default:
			// 	mt, message, err := c.ReadMessage()
			// 	if err != nil {
			// 		log.Println("read:", err)
			// 		break
			// 	}
			// 	log.Printf("recv: %s", message)
			// 	err = c.WriteMessage(mt, []byte("server"))
			// 	if err != nil {
			// 		log.Println("write:", err)
			// 		break
			// 	}

		}

	}

}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print( evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
</td></tr></table>
<div id="output"></div>
</body>
</html>
`))
