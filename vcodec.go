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

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	formType, r, err := riff.NewReader(bufio.NewReader(io.LimitReader(f, 65000)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RIFF(%s)\n", formType)
	if err := dump(r, ".\t"); err != nil {
		log.Fatal(err)
	}
}

func dump(r *riff.Reader, indent string) error {
	for {
		chunkID, chunkLen, chunkData, err := r.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if chunkID == riff.LIST {
			listType, list, err := riff.NewListReader(chunkLen, chunkData)
			if err != nil {
				return err
			}
			fmt.Printf("%sLIST(%s)\n", indent, listType)
			if err := dump(list, indent+".\t"); err != nil {
				return err
			}
			continue
		}
		b, err := ioutil.ReadAll(chunkData)
		if err != nil {
			return err
		}
		if string(b[0:4]) == "vids" {
			log.Printf("_%s_", string(b[4:8]))
		}
		fmt.Printf("%s%s %q\n", indent, chunkID, b)
	}
}
