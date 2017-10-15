package main

import (
    "net"
    "fmt"
    "strings"
    "os"
)
func chatMessage(rRef, name, mesg string) (string) {
	return fmt.Sprintf("CHAT: {}\nCLIENT_NAME: {}\nMESSAGE: {}\n\n")
}
type Client struct {
	conn		net.Conn
	addr		string
	uid			string
	name		string

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
	users		map[string][]Client
}


func main() {

    fmt.Println("Launching server...")

    // listen on all interfaces
    ln, _ := net.Listen("tcp", ":8080")

    // accept connection on port
    conn, _ := ln.Accept()
    var client = Client {conn,"testaddr","testuid","testname"}
    fmt.Print("client struc "+ client.name)
    // run loop forever (or until ctrl-c)
    for {
        // will listen for message to process ending in newline (\n)
        var buf = make([]byte, 1024)

        mesgLen, err := conn.Read(buf)


        if mesgLen ==0 {
            fmt.Printf("Connection closed by remote host\n")
        }
        checkError(err)
        // output message received
        fmt.Print("Message Received:", string(buf))
        // sample process for string received
        newmessage := strings.ToUpper(string(buf))
        // send new string back to client
        conn.Write([]byte(newmessage + "\n"))
    }

  

}
func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}