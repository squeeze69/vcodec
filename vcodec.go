// vcodec by Squeeze69
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/image/riff"
)

var exitValue = 0

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		fmt.Print("Usage: vcodec riff_file codec1 codec2 ...\ncodec name is case insensitive\n")
		os.Exit(2)
	}
	//A little trick to handle the "panic", os.Exit doesn't call deferred functions
	main2()
	fmt.Printf("Exit Value: %d\n", exitValue)
	os.Exit(exitValue)
}

func main2() {
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Codec: %s\n", r)
			codec := fmt.Sprintf("%s", r)
			for _, v := range flag.Args()[1:] {
				if strings.EqualFold(codec, v) {
					exitValue = 1
					break
				}
			}
		}
	}()
	formType, r, err := riff.NewReader(bufio.NewReader(io.LimitReader(f, 8192)))
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
