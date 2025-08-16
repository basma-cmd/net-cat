package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"tcp-chat/logo"
)

// Structure pour chaque client
type Client struct {
	Name string
	Conn net.Conn
}

var (
    clients   []Client   // slice pour stocker les clients
    clientsMu sync.Mutex // mutex pour protéger la slice
)

func main() {
	// 1) Écouter sur le port 8989
	ln, err := net.Listen("tcp", ":8989")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close() // fermer le server
	fmt.Println("Serveur TCP en écoute sur le port 8989")

	// 2) Accepter les clients dans une boucle
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Erreur accept:", err)
			continue
		}

		go handleClient(conn) // gérer chaque client dans une goroutine
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	_, _ = fmt.Fprint(conn, logo.Logo())

	_, _ = fmt.Fprint(conn, "[ENTER YOUR NAME]: ")
	reader := bufio.NewReader(conn) //stocke the name in tempo
	name, _ := reader.ReadString('\n') // read it jusqu'a \n 
	name = name[:len(name)-1] // retirer le \n
	if name == "" {
		name = "Anonymous"
	}

	// Ajouter le client à la liste
	clientsMu.Lock()
	clients = append(clients, Client{Name: name, Conn: conn})
	clientsMu.Unlock()

	fmt.Println(name, "vient de se connecter !")

	// Annoncer aux autres clients
	broadcast(fmt.Sprintf("%s a rejoint le chat...\n", name), conn)

	// Lire les messages
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break // le client s'est déconnecté
		}
		msg = msg[:len(msg)-1] // retirer le \n
		if msg != "" {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			broadcast(fmt.Sprintf("[%s][%s]: %s\n", timestamp, name, msg), nil)
		}
	}

	// Supprimer le client à la déconnexion
	clientsMu.Lock()
	for i, c := range clients {
    if c.Conn == conn { // trouver le client dans la slice
        clients = append(clients[:i], clients[i+1:]...) // supprimer
        break
    }
}
	clientsMu.Unlock()
	fmt.Println(name, "a quitté le chat.")
	broadcast(fmt.Sprintf("%s a quitté le chat...\n", name), conn)
}

// Fonction pour envoyer un message à tous les clients (sauf éventuellement l'envoyeur)
func broadcast(message string, ignoreConn net.Conn) {
    clientsMu.Lock()
    defer clientsMu.Unlock()

    for _, c := range clients {
        if c.Conn != ignoreConn {
            fmt.Fprint(c.Conn, message)
        }
    }
}

