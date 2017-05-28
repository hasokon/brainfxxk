package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	brainfxxk "github.com/hasokon/brainfxxk/lib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage command <filename>\n")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer file.Close()

	bf := brainfxxk.New()

	reader := bufio.NewReader(file)
	for {
		input, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		for _, c := range input {
			err = bf.Add(c)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		}
	}

	bf.Run()
}
