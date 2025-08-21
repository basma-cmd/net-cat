package main

import (
	"fmt"
	"net"
	"time"

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
	resp, err := net.Listen("tcp", ":8989")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Close()

	fmt.Println("server listen in port 8989")
	for {
		conn, err := resp.Accept()
		if err != nil {
			fmt.Println("erro:", err)
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
    joinMsg := name + " has joined our chat...\n"
    broadcast(joinMsg, memberId)

    // boucle pour écouter les messages de ce client
    for {
        buff := make([]byte, 1024)
        n, err := conn.Read(buff)
        if err != nil {
            fmt.Println("error:", err)
            return
        }

        message := strings.TrimSpace(string(buff[:n]))
        if message == "" {
            continue
        }

        now := time.Now().Format("2006-01-02 15:04:05")
        fullMsg := fmt.Sprintf("[%s][%s]: %s\n", now, name, message)

        // envoyer à tous sauf lui-même
        broadcast(fullMsg, memberId)
    }
}

// broadcast à tous sauf l’expéditeur
func broadcast(msg string, senderId int) {
    for _, m := range members {
        if m.Id != senderId {
            _, err := m.Conn.Write([]byte(msg))
            if err != nil {
                fmt.Println("Error writing to", m.Name, ":", err)
            }
        }
    }
}
