package main

import "net"
import "fmt"
import "bufio"
import "os" 
func main() {

  // connect to this socket
  conn, _ := net.Dial("tcp", "127.0.0.1:8080")
  for { 
    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')
     _ = text

    var c =0
    for {
        var buf = make([]byte,1024)
        c +=1
         
        // send to socket
        fmt.Fprintf(conn,"JOIN_CHATROOM: chat3\nCLIENT_IP: 0\nPORT: 0\nCLIENT_NAME: client2\n")


        //fmt.Fprintf(conn,"LEAVE_CHATROOM: chat3\nJOIN_ID: 0\nCLIENT_NAME: client2\n")
        // listen for reply
        conn.Read(buf)
    
        fmt.Print("Message from server: "+string(buf))
        break
    }
    
  }
}