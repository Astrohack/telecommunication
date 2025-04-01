package main

import (
	"fmt"
	"math/bits"
	"net"
)

var H = []uint{
	0b1111000010000000,
	0b1100110001000000,
	0b1010101000100000,
	0b0101011000010000,
	0b1110100100001000,
	0b1001010100000100,
	0b0111101100000010,
	0b1110011100000001,
}

var lookupH = []uint{
	0b10000000,
	0b01000000,
	0b00100000,
	0b00010000,
	0b00001000,
	0b00000100,
	0b00000010,
	0b00000001,
	0b11110000,
	0b11001100,
	0b10101010,
	0b01010110,
	0b11101001,
	0b10010101,
	0b01111011,
	0b11100111,
}

// makeLookupTable makes lookup table to quickly find error positions by syndrome
func makeLookupTable(G []uint) map[byte]uint16 {
	rows := len(G)
	lookupTable := make(map[byte]uint16)
	for i, val := range G {
		lookupTable[byte(val)] = 1 << i
		for j := i + 1; j < rows; j++ {
			key := byte(val ^ G[j])              // syndrome as key
			value := uint16((1 << i) | (1 << j)) // error positions as value
			lookupTable[key] = value
		}
	}
	return lookupTable
}

// parity returns quantity of ones in given value
func parity(word uint16) byte {
	return byte(bits.OnesCount(uint(word)) % 2)
}

// bitDecode returns syndrome for given codeword and H matrix
func bitDecode(H []uint, word uint16) byte {
	n := len(H)
	var syndrome byte = 0
	for idx := range n {
		syndrome |= parity(uint16(H[idx])&word) << idx
	}
	return syndrome
}

func decode(received uint16, correctionLookupTable map[byte]uint16) byte {
	syndrome := bitDecode(H, received)
	if syndrome != 0 {
		fmt.Printf("\nsyndrome [%b] (%b) (%b) -> errors at %b\n", syndrome, received>>8, byte((received^correctionLookupTable[syndrome])>>8), correctionLookupTable[syndrome])
	}
	return byte((received ^ correctionLookupTable[syndrome]) >> 8)
}

func main() {
	correctionLookupTable := makeLookupTable(lookupH)

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Błąd ResolveUDPAddr:", err)
		return
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Błąd DialUDP:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error while receiving message:", err)
		}

		for i := 0; i < n; i += 2 {
			if i+1 == n {
				fmt.Println("Dropping some bytes due to loss of data")
				continue
			}
			var received = uint16(buffer[i])<<8 | uint16(buffer[i+1])
			decoded := decode(received, correctionLookupTable)

			fmt.Printf("%c", decoded)
		}
	}
}
