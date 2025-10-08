package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func check(e error) {
	if e != nil {
		//panic(e)
		log.Fatalf("could not open %s\n", e)
	}
}

const inputFilePath = "messages.txt"

func main_from_file() {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	linesChan := getLinesChannel(f)

	for line := range linesChan {
		fmt.Println("read:", line)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not obtain listener: %s\n", err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %s\n", err)
		}
		fmt.Printf("Connection accepted \n")

		linesChan := getLinesChannel(connection)

		for line := range linesChan {
			fmt.Println("read:", line)
		}
	}

}
