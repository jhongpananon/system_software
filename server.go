package main

import (
    "bufio"
    "fmt"
    "net"
)

var client_map map[*client_S]int

type client_S struct {
    outgoing   chan string
    rx_io      *bufio.Reader
    tx_io      *bufio.Writer
    conn       net.Conn
    client_ptr *client_S
}

func (client *client_S) Read() {
    for {
        line, err := client.rx_io.ReadString('\n')
        if err == nil {
            if client.client_ptr != nil {
                client.client_ptr.outgoing <- line
            }
            fmt.Println(line)
        } else {
            break
        }

    }

    client.conn.Close()
    delete(client_map, client)
    if client.client_ptr != nil {
        client.client_ptr.client_ptr = nil
    }
    client = nil
}

func (client *client_S) Write() {
    for data := range client.outgoing {
        client.tx_io.WriteString(data)
        client.tx_io.Flush()
    }
}

//func (client *client_S) Listen() {
    // Start the two goroutines to read/write
//    go client.Read()
//    go client.Write()
//}

func Newclient(client_ptr net.Conn) *client_S {
    tx_io := bufio.NewWriter(client_ptr)
    rx_io := bufio.NewReader(client_ptr)

    client := &client_S {
        outgoing: make(chan string),
        conn:     client_ptr,
        rx_io:    rx_io,
        tx_io:    tx_io,
    }
    //client.Listen()

    return client
}

func main() {
    client_map = make(map[*client_S]int)

    // net Listen for a TCP socket on this port
    listener, _ := net.Listen("tcp", ":4005")
    for {
        // Block on accept until client connects
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(err.Error())
        }

        // Start a new client connection
        client := Newclient(conn)

        // Loop through the client_map list
        for clientList, _ := range client_map {
            if clientList.client_ptr == nil {
                client.client_ptr = clientList
                clientList.client_ptr = client
                fmt.Println("Connected")
            }
        }
        client_map[client] = 1
        fmt.Println(len(client_map))
    }
}