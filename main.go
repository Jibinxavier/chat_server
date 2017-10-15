package main

import (
    "net"
    "fmt"
   // "strings"
    "os"
)
var serverIP    string
var serverPort  string

func chatMessage(rRef, name, mesg string) (string) {
	return fmt.Sprintf("CHAT: {}\nCLIENT_NAME: {}\nMESSAGE: {}\n\n")
}
type Client struct {
	conn		net.Conn
	addr		string
	uid			string
    name		string
    incoming    chan string
	outgoing    chan string

}
type Mesg struct {
    mesgFrom 	string
	chatId 		string
	clientName 	string
	mesg 		string
}

type ChatRoom struct {
	id			string
	name 		string
	users		map[string][]Client // key would be the name of the client
}
 
func (client *Client) Read() {
    var buf = make([]byte, 1024)
    
    for {  
        mesgLen, err := client.conn.Read(buf) 

        if mesgLen ==0 {
            fmt.Printf("Connection closed by remote host\n")
        }
        checkError(err)
        // output message received
        fmt.Print("Message Received:", string(buf)) 
    }
    client.incoming <- string(buf)
	 
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.conn.Write([]byte(data + "\n"))
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn) *Client {  
	client := &Client{ 
		conn:       connection,
        addr:		"testAddr",
        uid:		"testAddr",
        name:		"testname",
        incoming:   make(chan string),
		outgoing:   make(chan string),
	}

	client.Listen()

	return client
}

// func parseMesg(mesg string){
//     var data = strings.Split(mesg, "\n")
//     var resp string
//     if strings.Contains(data[0], "HELO") {
//         resp := resp + fmt.Sprintf("IP:%s\nPort:%s\nStudentID:13321596\n", serverIP, serverPort)
//     }
//     fmt.Println("Launching server..." + resp)

// }
func clientHandle(conn net.Conn){
    

}
 

func main() {
    serverIP   := "127.0.0.1"
    serverPort := "8080" 
    fmt.Println("Launching server...")

    // listen on all interfaces
     
    
     
    ln, _ := net.Listen("tcp", serverIP + ":" + serverPort)
    
    // accept connection on port
    
    //var client = Client {conn,"testaddr","testuid","testname"}
    

    for {
		conn, _ := ln.Accept()
		NewClient(conn)
	}
  

}
func checkError(err error) {


    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}