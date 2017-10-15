package main

import "net"
import "fmt"
//import "bufio"
import "strings" // only needed below for sample processing
//import "io/ioutil"
import "os"
func main() {

  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", ":8080")

  // accept connection on port
  conn, _ := ln.Accept()

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