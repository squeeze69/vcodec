// vcodec by Squeeze69
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/image/riff"
)

var exitValue = 0

func main() {
	if len(os.Args) < 3 {
		fmt.Print("Usage: vcodec riff_file codec1 codec2 ...\ncodec name case DO matter\n")
		os.Exit(2)
	}
	main2()
	fmt.Printf("Exitvalue: %d\n", exitValue)
	os.Exit(exitValue)
}

func main2() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Codec: %s\n", r)
			for _, v := range os.Args[2:] {
				if r == v {
					exitValue = 1
					break
				}
			}
		}
	}()
	formType, r, err := riff.NewReader(bufio.NewReader(io.LimitReader(f, 16384)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RIFF(%s)\n", formType)
	if err := scanriff(r); err != nil {
		log.Fatal(err)
	}
}

func scanriff(r *riff.Reader) error {
	for {
		chunkID, chunkLen, chunkData, err := r.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if chunkID == riff.LIST {
			_, list, err := riff.NewListReader(chunkLen, chunkData)
			if err != nil {
				return err
			}
			if err := scanriff(list); err != nil {
				return err
			}
			continue
		}
		b, err := ioutil.ReadAll(chunkData)
		if err != nil {
			return err
		}
		if string(b[0:4]) == "vids" {
			panic(string(b[4:8]))
		}
	}
}
