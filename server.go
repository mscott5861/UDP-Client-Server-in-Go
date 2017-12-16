package main

import (
    "fmt"
    "net"
    "os"
    "encoding/binary"
    "time"
    "strconv"
)

func ParseErrorResponse(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(0)
    }
}

func main() {
    ServerPort := ":" + os.Args[1]
    ServerAddr, err := net.ResolveUDPAddr("udp", ServerPort)
    ParseErrorResponse(err)

    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    ParseErrorResponse(err)
    defer ServerConn.Close()

    // We'll allocate enough for the default UDP MTU limit, though
    // we won't need it.
    buf := make([]byte, 65535)
    timesMsgReceived := 0

    for {
        n, _, err := ServerConn.ReadFromUDP(buf)
        timesMsgReceived++
        // Our message occupies bytes 18 -- length of message in the payload
        fmt.Println("Received: ", string(buf[18:n]))

        // Handle parsing our timestamp. An 8-byte Unix timestamp was sent at
        // bytes 4-12
        timestamp := int64(binary.BigEndian.Uint64(buf[4:12]))
        clientTimestamp := time.Unix(timestamp, 0)

        // Now get the server's own Unix timestamp
        serverTimestamp := time.Unix(time.Now().Unix(), 0)
        differenceInTime := serverTimestamp.Sub(clientTimestamp)
        fmt.Println("Time elapsed between tx/rx: ", differenceInTime)

        // Get the IP address stored at bytes 12-16
        ip := binary.BigEndian.Uint32(buf[12:16])
        port := binary.BigEndian.Uint16(buf[16:18])
        portStr := strconv.Itoa(int(port))
        ipB := make(net.IP, 4)
        binary.BigEndian.PutUint32(ipB, ip)
        fmt.Println("IPv4 of sender: ", net.IP.String(ipB) + ":" + portStr)
        fmt.Println("Times message received: ", timesMsgReceived)
        fmt.Println("\n")

        if err != nil {
            fmt.Println("Error: ", err)
        }
    }
}




