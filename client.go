package main

import (
    "bufio"
    "fmt"
    "net"
)

var client_map map[*client_S]int

type client_S struct {
    // incoming chan string
    outgoing   chan string
    rx_io      *bufio.Reader
    tx_io      *bufio.Writer
    conn       net.Conn
    connection *client_S
}

func (client *client_S) Read() {
    for {
        line, err := client.rx_io.ReadString('\n')
        if err == nil {
            if client.connection != nil {
                //client.connection.outgoing <- line
            }
            fmt.Println(line)
        } else {
            break
        }

    }

    client.conn.Close()
    delete(client_map, client)
    if client.connection != nil {
        client.connection.connection = nil
    }
    client = nil
}

func (client *client_S) Write() {
    for data := range client.outgoing {
        client.tx_io.WriteString(data)
        client.tx_io.Flush()
    }
}

func (client *client_S) Listen() {
    go client.Read()
    go client.Write()
}

func NewClient(connection net.Conn) *client_S {
    tx_io := bufio.NewWriter(connection)
    rx_io := bufio.NewReader(connection)

    client := &client_S {
        // incoming: make(chan string),
        outgoing: make(chan string),
        conn:     connection,
        rx_io:    rx_io,
        tx_io:    tx_io,
    }
    client.Listen()

    return client
}

func main() {
    client_map = make(map[*client_S]int)
    for {
        // net TCP dial to the remote server
        conn, err := net.Dial("tcp", "squad:4005")
        if err != nil {
            fmt.Println(err.Error())
        }

        // New client connected
        client := NewClient(conn)
        for clientList, _ := range client_map {
            if clientList.connection == nil {
                client.connection = clientList
                clientList.connection = client
                fmt.Println("Connected")
            }
        }
        client_map[client] = 1
        fmt.Println(len(client_map))
    }
}