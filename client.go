package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"math/rand/v2"
	"net"
	"os"
)

var G = []byte{
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

// parity returns quantity of ones in given value
func parity(word uint16) uint16 {
	return uint16(bits.OnesCount(uint(word)) % 2)
}

// bitEncode returns codeword for given word and matrix generator
func bitEncode(G []byte, word byte) (byte, byte) {
	n := len(G)
	var codeword uint16 = 0
	for idx := range n {
		codeword = codeword | (parity(uint16(G[n-1-idx]&word)) << idx)
	}
	return byte(codeword >> 8), byte(codeword & 0xFF)
}

func encode(word string) []byte {
	wordLen := len(word)
	encoded := make([]byte, wordLen*2)
	for i, char := range word {
		encoded[i*2], encoded[i*2+1] = bitEncode(G, byte(char))
		if char > 14 {
			fmt.Printf("; sending char (%c) %d: %b %b\n", char, i, encoded[i*2], encoded[i*2+1])
		}
	}
	return encoded
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Błąd ResolveUDPAddr:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Błąd DialUDP:", err)
		return
	}
	defer conn.Close()

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		encoded := encode(text)
		for idx, char := range encoded {
			encoded[idx] = char ^ (1 << rand.UintN(8))
		}
		_, err := conn.Write(encoded)
		if err != nil {
			fmt.Println("Error while sending message:", err)
		}
	}
}
