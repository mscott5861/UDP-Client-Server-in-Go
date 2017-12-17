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

    argCount := len(os.Args[1:])
    if (argCount != 1) {
        fmt.Print("\nThe server expects exactly one command line parameter:\n\n\t")
        fmt.Print("1. The port for the server to listen on\n\nExample invocation:\n\n\t")
        fmt.Println("go run server.go 9887\n")
        os.Exit(1)
    }

    // Gather the port we'll be listen on, check that it's in good shape, and start eavesdropping.
    ServerPort := ":" + os.Args[1]
    ServerAddr, err := net.ResolveUDPAddr("udp", ServerPort)
    ParseErrorResponse(err)
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    ParseErrorResponse(err)

    // We'll be charitable and let people know we're eavesdropping
    fmt.Println("\n\nListening on port " + os.Args[1] + " (\"Ctrl+C\" to quit)\n\n")
    defer ServerConn.Close()

    // Allocate enough for the default UDP MTU limit, though we won't anything close to it.
    // TCPDump pegs these payload for these datagrams at 109 bytes in length.
    buf := make([]byte, 65535)
    timesMsgReceived := 0

    for {
        n, _, err := ServerConn.ReadFromUDP(buf)
        timesMsgReceived++
        //------------------------------------------------------------------
        // BYTES 18-n
        //------------------------------------------------------------------
        // Our message occupies bytes 18 -> length of message in the payload
        //------------------------------------------------------------------
        fmt.Println("Received: ", string(buf[18:n]))

        //------------------------------------------------------------------
        // BYTES 4-12
        //------------------------------------------------------------------
        // Handle parsing our timestamp. An 8-byte Unix timestamp was sent at
        // bytes 4-12
        //------------------------------------------------------------------
        timestamp := int64(binary.BigEndian.Uint64(buf[4:12]))
        clientTimestamp := time.Unix(timestamp, 0)

        // Now get the server's own Unix timestamp
        serverTimestamp := time.Unix(time.Now().Unix(), 0)
        // Subtract the difference. On the LAN (or loopback device), this difference will
        // return 0, as the tx/rx time will be negligible and precision of the Unix
        // timestamp is at the seconds level (not milli- or mico)
        differenceInTime := serverTimestamp.Sub(clientTimestamp)
        fmt.Println("Time elapsed between tx/rx: ", differenceInTime)

        //------------------------------------------------------------------
        // BYTES 12-16
        //------------------------------------------------------------------
        // Get the IP address stored at bytes 12-16
        //------------------------------------------------------------------
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




