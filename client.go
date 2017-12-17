package main

import (
    "os"
    "fmt"
    "net"
    "time"
    "encoding/binary"
    "strconv"
    "hash/adler32"
    "bytes"
)

func ParseErrorResponse(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
    }
}

// We take an int large enough to accommodate the largest we'll pass (uint64),
// then allocate an appropriate sized byte array, typecast the argument, and
// return.
func GetIntBigEndianBytes(numBytes int, arg uint64) (buf []byte) {
    buf = make([]byte, numBytes)
    if numBytes == 2 {
        binary.BigEndian.PutUint16(buf, uint16(arg))
    } else if numBytes == 4 {
        binary.BigEndian.PutUint32(buf, uint32(arg))
    } else if numBytes == 8 {
        binary.BigEndian.PutUint64(buf, uint64(arg))
    }
    return
}

func GetStrBigEndianBytes(str string) (buf []byte) {
    msg := []byte(str)
    msgBuf := new(bytes.Buffer)
    err := binary.Write(msgBuf, binary.BigEndian, msg)
    ParseErrorResponse(err)
    buf = msgBuf.Bytes()
    return
}

// Take the IP address as a net.IP object, create a byte array buffer,
// convert the net.IP object to a 4-byte representation, then return
// a Big Endian representation of the address

func GetIPBigEndianBytes(ipAddr net.IP) (buf []byte) {
    buf = make([]byte, 4)
    ip := binary.BigEndian.Uint32(net.IP.To4(ipAddr))
    binary.BigEndian.PutUint32(buf, ip)
    return
}

func main() {
    // First, let's gather our own local IP Address. The address returned by InterfaceAddrs
    // at index 0 will always be the loopback device. We're looking for the IP on the LAN,
    // at index 1.
    myIPCIDR, err := net.InterfaceAddrs()
    ParseErrorResponse(err)
    myIP,_,err := net.ParseCIDR(myIPCIDR[1].String())
    ParseErrorResponse(err)

    argCount := len(os.Args[1:])
    if (argCount != 2) {
        fmt.Print("\nThe client expects exactly two command line parameters:\n\n\t")
        fmt.Print("1. The port for the client to listen on\n\t2. The sever's IP:Port.")
        fmt.Println("\n\nExample invocation:\n\n\tgo run client.go 9999 192.168.1.10:9887\n")
        os.Exit(1)
    }

    // Get the client port passed as the first argument, then concatenate it
    // to the the client's LAN address
    clientPort := os.Args[1]
    clientIP := myIP.String() + ":" + clientPort
    // Get the server IP address passed as the second argument to the program
    serverIP := os.Args[2]


    ServerAddr, err := net.ResolveUDPAddr("udp", serverIP)
    ParseErrorResponse(err)

    LocalAddr, err := net.ResolveUDPAddr("udp", clientIP)
    ParseErrorResponse(err)

    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    ParseErrorResponse(err)

    defer Conn.Close()
    for {
        // Handle creating our message
        msg := "\"It always seems impossible until it's done.\" - Nelson Mandella"
        msgBuf := GetStrBigEndianBytes(msg)

        // Handle pushing our port
        portInt, _ := strconv.Atoi(clientPort)
        portBuf := GetIntBigEndianBytes(2, uint64(portInt))
        concat1 := append(portBuf[:], msgBuf[:]...)

        // Handle pushing our IP Address. Feed it a String copy of the IP address we
        // gathered earlier (our IP on the LAN).
        ipBuf := GetIPBigEndianBytes(net.ParseIP(myIP.String()))
        concat2 := append(ipBuf[:], concat1[:]...)

        // Handle creating our timestamp. Not much to see here.
        timeBuf := GetIntBigEndianBytes(8, uint64(time.Now().Unix()))
        concat3 := append(timeBuf[:], concat2[:]...)

        // Handle creating our Adler32 checksum. It is passed the current payload,
        // which it calculates a hash against. The value it returns is converted 
        // into a Big Endian byte array, and then concatenated with the payload. 
        checksumBuf := GetIntBigEndianBytes(4, uint64(adler32.Checksum(concat3)))
        concat4 := append(checksumBuf[:], concat3[:]...)

        _, err2 := Conn.Write(concat4)
        ParseErrorResponse(err2)

        // Send our packet every two seconds
        time.Sleep(time.Second * 2)
    }
}

