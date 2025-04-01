package main

import (
	"fmt"
	"math/bits"
)

var bitH = []uint{
	0b1111000010000000,
	0b1100110001000000,
	0b1010101000100000,
	0b0101011000010000,
	0b1110100100001000,
	0b1001010100000100,
	0b0111101100000010,
	0b1110011100000001,
}

var bitLookupH = []uint{
	0b00001111,
	0b00110011,
	0b01010101,
	0b01101010,
	0b10010110,
	0b10101011,
	0b11011011,
	0b11101101,
}

var bitG = []uint{
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
func makeLookupTable(G []uint) map[uint]uint {
	rows := len(G)
	lookupTable := make(map[uint]uint)
	for i, val := range G {
		lookupTable[val] = 1 << i
		for j := i + 1; j < rows; j++ {
			key := val ^ G[j]                  // syndrome as key
			value := uint((1 << i) | (1 << j)) // error positions as value
			lookupTable[key] = value
		}
	}
	return lookupTable
}

// parity returns quantity of ones in given value
func parity(word1 uint) uint {
	return uint(bits.OnesCount(word1) % 2)
}

// bitEncode returns codeword for given word and matrix generator
func bitEncode(G []uint, word uint) uint {
	n := len(G)
	var codeword uint = 0
	for idx := range n {
		codeword = codeword | (parity(G[n-1-idx]&word) << idx)
	}
	return codeword
}

// bitDecode returns syndrome for given codeword and H matrix
func bitDecode(H []uint, word uint) uint {
	n := len(H)
	var syndrome uint = 0
	for idx := range n {
		syndrome |= parity(H[idx]&word) << idx
	}
	return syndrome
}

func main() {
	correctionLookupTable := makeLookupTable(bitG)

	for key, val := range correctionLookupTable {
		fmt.Printf("%b -> %b\n", key, val)
	}

	bitMessage := uint(0b00110111)

	encoded := bitEncode(bitG, bitMessage)
	fmt.Printf("\nEncoded message %b: %b\n", bitMessage, encoded)

	//received := encoded ^ 0b00001000000010 // introduce error
	received := encoded ^ 0b00011000000000 // introduce error

	syndrome := bitDecode(bitH, received)
	fmt.Printf("\nReceived %b, syndrome is: %b, errors are at positions: %b\n", received, syndrome, correctionLookupTable[syndrome])

	correctedMessage := (received >> 8) ^ correctionLookupTable[syndrome]
	fmt.Printf("\nCorrected message: %b\n", correctedMessage)
}
