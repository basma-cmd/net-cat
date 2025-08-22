package main

import (
	"fmt"
	"net"
	"time"

	//"time"

	//"os"
	"strings"
	//"sync"
)

var s = `Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML`

type Members struct {
	Id   int
	Name string
	Conn net.Conn
}

var (
	latestID int
	// mu       sync.Mutex
	members []Members
)

func main() {
	resp, err := net.Listen("tcp", ":8987")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Close()

	fmt.Println("server listen in port 8987")
	for {
		conn, err := resp.Accept()
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("le client de ip:", conn.RemoteAddr(), ", a connecter")

	// envoi du logo + prompt
	_, err := conn.Write([]byte(s + "\n" + "[Enter your name]: "))
	if err != nil {
		fmt.Println("error", err)
		return
	}

	// lecture du nom
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	name := strings.TrimSpace(string(buf[:n]))

	// création du membre
	memberId := latestID
	latestID++
	newMember := Members{Id: memberId, Name: name, Conn: conn}
	members = append(members, newMember)

	// informer les autres
	joinMsg := "\n" + name + " has joined our chat...\n"

	broadcast(conn, joinMsg, memberId, name)

	// boucle pour écouter les messages de ce client
	for {
		now := time.Now()
		formatted := now.Format("2006-01-02 15:04:05")
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("error", err)
		}
		message := strings.TrimSpace(string(buff[:n]))
		message = "\n" + "[" + formatted + "]" + "[" + name + "]:" + message + "\n"
		broadcast(conn, message, memberId, name)
	}
}

// broadcast à tous sauf l’expéditeur
func broadcast(conn net.Conn, msg string, senderId int, name string) {
	for _, m := range members {
		if m.Id != senderId {
			_, err := m.Conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error writing to", m.Name, ":", err)
				return
			}
			
		}
		now := time.Now()
		formatted := now.Format("2006-01-02 15:04:05")
		message := "\n" + "[" + formatted + "]" + "[" + m.Name + "]:"
		//fmt.Println(message)
		_, err := m.Conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to", name, ":", err)
			return
	
		}
	}
}
