package main

import (
    "net"
    "fmt"
    "strings"
    "os"
    "sync"
)
var serverIP    string
var serverPort  string

func chatMessage(roomRef, clientName, mesg string) (string) {
    return fmt.Sprintf("CHAT: %s\nCLIENT_NAME: %s\nMESSAGE: %s\n\n",
                                                            roomRef,
                                                            clientName, 
                                                            mesg)
}
func errorMessage(errorNum int, mesg string) string {
    return fmt.Sprintf("ERROR_CODE:%d\nERROR_DESCRIPTION: %s\n",errorNum, mesg)
}
type Session struct {
    mu          sync.Mutex
    chatRooms   map[string]*ChatRoom // key name of chat room
}
type Client struct {
	conn		net.Conn
	addr		string
	uid			string
    name		string
    incoming    chan string
    outgoing    chan string
    sess        *Session

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
    clients		map[string]*Client // key would be the name of the client
    incoming    chan string
	outgoing    chan string
}
func NewChatRoom(id string, name string ) *ChatRoom {  
	chatroom := &ChatRoom{ 
		id:			id,
        name: 		name,
        clients:	make(map[string]*Client) , 
        incoming:   make(chan string),
		outgoing:   make(chan string),
    } 
    
	return chatroom
}
 
 // create a thread for each chatroom
//
func (chatRoom *ChatRoom) Broadcast(data string) {
	for _, client := range chatRoom.clients {
		client.outgoing <- data
	}
}
func (sess *Session) getChatroom( roomName string) (*ChatRoom,bool){

    for _, room := range sess.chatRooms {
        if  room.name == roomName {
             return room, true
        } 
    }
    return nil, false
}
func (sess *Session) getChatrooWRef( ref string) (*ChatRoom,bool){
    
        for _, room := range sess.chatRooms {
            if  room.id == ref {
                 return room, true
            } 
        }
        return nil, false
    }
func (client *Client) chat(mesg string, clientName string, roomRef string, joinId string){
    client.sess.mu.Lock()
    defer client.sess.mu.Unlock()

    
    if client.name == clientName {
        chatroom := client.sess.chatRooms[roomRef]
        broadcastMesg :=chatMessage(roomRef,client.name,mesg)
        chatroom.Broadcast(broadcastMesg)
        //close channels 
    } else {
        client.outgoing <- errorMessage( 24, "User name not found")
    }

}

func (client *Client) joinChatroom(roomName string, userName string )  { 
    client.sess.mu.Lock()
    defer client.sess.mu.Unlock()

    var clientMesg      string
    var broadcastMesg   string
    var chatroom        *ChatRoom 
    var roomRef         string = "0"

    chatroom, ok :=  client.sess.getChatroom( roomName)
    // create the chat room if it is new
    if  ok == false {
        //fmt.Printf("New chatroom \n")
        roomRef := fmt.Sprint(len(client.sess.chatRooms))  
        fmt.Printf("New chatroom %d\n",len(client.sess.chatRooms))
        newChatRoom := NewChatRoom(roomRef, roomName )
        client.sess.chatRooms[roomRef] = newChatRoom // add values to map 
        chatroom = newChatRoom
    }   
    // add clients if new to the chat room 
    _, ok = chatroom.clients[userName]
    if ok {
        clientMesg = " you already are in the chat room "  

    }else {  
        chatroom.clients[userName] = client // add client to the chatroom
        clientMesg = fmt.Sprintf("JOINED_CHATROOM: %s\nSERVER_IP: %s\nPORT:%s\nROOM_REF:%s\nJOIN_ID: %s",chatroom.name,
                                                                                                    serverIP,
                                                                                                    serverPort,
                                                                                                    chatroom.id,
                                                                                                  client.uid ) 
        broadcastMesg = chatMessage(roomRef,client.name,fmt.Sprintf("client %s has joined this chatroom.",userName))
    }
    client.outgoing <- clientMesg   // send client notification
    chatroom.Broadcast(broadcastMesg) // notification to the whole chat room
}
func (client *Client) disconnect(clientName string){
    if client.name == clientName {
        client.conn.Close()
        client.sess.mu.Unlock()
        //close channels 
    } else {
        client.outgoing <- errorMessage( 24, "User name not found")
    }
}
func (client *Client) leaveChatroom(roomRef string, joinId string, userName string)  {

    client.sess.mu.Lock()
    defer client.sess.mu.Unlock()

    var clientMesg      string
    var broadcastMesg   string 
    // joinid and uid are the same
    chatroom, ok :=  client.sess.getChatrooWRef(roomRef)
    //chatroom, ok :=  client.sess.chatRooms[roomRef]
    // can not find chat
    if ok == false {
        fmt.Printf(" Chat room not found%s",roomRef)
        clientMesg = errorMessage(1, "Unknown chat room")
         
    }else {
        _, ok = chatroom.clients[userName]; 
        if ok {
            delete(chatroom.clients, userName) 


            clientMesg = fmt.Sprintf("LEFT_CHATROOM: %s\nJOIN_ID: %s",chatroom.name, joinId) 

            fmt.Printf(clientMesg)
            broadcastMesg = chatMessage(roomRef,client.name,fmt.Sprint("client %s has left this chatroom.",client.name))
            _ = broadcastMesg
            chatroom.Broadcast(broadcastMesg) // notification to the whole chat room
        }else {  
            clientMesg = errorMessage( 24, "User name not found") 
            // can not find user name
        } 
    }
    
    client.outgoing <- clientMesg   // send client notification
    
}

func (client *Client) updateClient(name string){
    /*
        When the client is created, it does not have
        the name. However certain commands such as
        JOIN_CHATROOM have the clients name
    */ 
    if client.name == "defaultname" {
        client.name = name
    }
}


func (client *Client) parseMesg(mesg string ){
    var data [] string= strings.Split(mesg, "\n")
     
    var clientName  string 
    var roomName    string
    var roomRef     string 
    var joinId      string

    if strings.Contains(data[0], "HELO") {
        client.outgoing <- fmt.Sprintf("IP:%s\nPort:%s\nStudentID:13321596\n", serverIP, serverPort)
        
    } else if strings.Contains(data[0], "JOIN_CHATROOM") { 
        
        roomName    = strings.Trim(strings.Split(data[0], ":")[1], " ") // as per protocol structure 
        clientName  = strings.Trim(strings.Split(data[3], ":")[1], " ")

        client.updateClient(clientName)
        client.joinChatroom(roomName, clientName)

    } else if strings.Contains(data[0], "LEAVE_CHATROOM"){
        roomRef     = strings.Trim(strings.Split(data[0], ":")[1], " ")
        joinId      = strings.Trim(strings.Split(data[1], ":")[1], " ")
        clientName  = strings.Trim(strings.Split(data[2], ":")[1], " ")

        client.leaveChatroom(roomRef , joinId, clientName) 
    } else if strings.Contains(data[0], "DISCONNECT") {
        clientName  = strings.Trim(strings.Split(data[2], ":")[1], " ")

        client.disconnect(clientName)
    } else if strings.Contains(data[0], "CHAT"){
        roomRef     = strings.Trim(strings.Split(data[0], ":")[1], " ")
        joinId      = strings.Trim(strings.Split(data[1], ":")[1], " ")
        clientName  = strings.Trim(strings.Split(data[2], ":")[1], " ")
        mesg        = strings.Trim(strings.Split(data[3], ":")[1], " ")

        client.chat(mesg, clientName,roomRef,joinId)
    }else if  strings.Contains(data[0], "KILL_SERVICE"){
        os.Exit(1)
    }

}

func (client *Client) Read() {
    var buf = make([]byte, 1024)
    
    for {  
        mesgLen, err := client.conn.Read(buf) 

        if mesgLen ==0 {
            fmt.Printf("Connection closed by remote host\n")
            break
        }
        checkError(err)
        // output message received
        client.parseMesg(string(buf))
       
    }
   
	 
}

func (client *Client) Write() {
    fmt.Print("writing to channel") 
	for data := range client.outgoing {
		client.conn.Write([]byte(data + "\n"))
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn,uid int, sess *Session) *Client {  
	client := &Client{ 
		conn:       connection,
        addr:		"testAddr",
        uid:		fmt.Sprint(uid),
        name:		"defaultname",
        incoming:   make(chan string),
        outgoing:   make(chan string),
        sess:       sess,
	}

	client.Listen()

	return client
}

 

func main() {

    if (len(os.Args) !=3) {
        fmt.Println("ip port")
        os.Exit(1)
    }
    parms := os.Args[1:]
    serverIP   = parms[0]
    serverPort = parms[1] 
    fmt.Println("Launching server...")

    // listen on all interfaces
     
    currSession := &Session{chatRooms: make(map[string]*ChatRoom)}
     
    ln, _ := net.Listen("tcp", serverIP + ":" + serverPort)
    
    // accept connection on port 
    clientCount := 0

    for {
        conn, _ := ln.Accept() // waits for new connection
        NewClient(conn,clientCount,currSession)
        clientCount +=1
	}
  

}
func checkError(err error) {


    if err != nil {
        //fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        //os.Exit(1)
    }
}