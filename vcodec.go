// vcodec by Squeeze69
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/image/riff"
)

// biCompression is the fourCC code
/*
type tagBITMAPINFOHEADER struct {
	biSize          uint32
	biWidth         int32
	biHeight        int32
	biPlanes        uint16
	biBitCount      uint16
	biCompression   uint32
	biSizeImage     uint16
	biXPelsPerMeter int32
	biYPelsPerMeter int32
	biClrUsed       uint32
	biClrImportant  uint32
}
*/

// value for os.Exit
var exitValue = 0

// list of codecs to scan for
var codecList []string

// file to read codecList from (optional)
var codecListfile string

func main() {
	flag.StringVar(&codecListfile, "c", "", "-c codeclistfile (one per line)")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Print("Usage: vcodec [options] riff_file [optional codec1 codec2 ...]\ncodec name is case insensitive\n")
		flag.Usage()
		os.Exit(2)
	}
	//the next step is to allow to load the codecList from a file, as an alternative
	if flag.NArg() > 1 {
		codecList = flag.Args()[1:]
	} else if codecListfile != "" {
		buf, err := os.ReadFile(codecListfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		scanner := bufio.NewScanner(strings.NewReader(string(buf)))
		for scanner.Scan() {
			codecList = append(codecList, scanner.Text())
		}
	}

	//A little trick to handle the "panic", os.Exit doesn't call deferred functions
	main2()
	fmt.Printf("Exit Value: %d\n", exitValue)
	os.Exit(exitValue)
}

// Second main, to overcame os.Exit behaviour (deferred func execution)
func main2() {
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// check for a matching codec, case insensitive
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Codec: %s\n", r)
			codec := fmt.Sprintf("%s", r)
			for i := range codecList {
				if strings.EqualFold(codec, codecList[i]) {
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

// scan riff for chunk data
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
		b, err := io.ReadAll(chunkData)
		if err != nil {
			return err
		}
		if string(b[0:4]) == "vids" {
			codec := string(b[4:8])
			chunkID, _, chunkData, err := r.Next()
			if err != nil {
				panic(codec)
			}
			if fmt.Sprintf("%s", chunkID) == "strf" {
				b, err = io.ReadAll(chunkData)
				if err == nil {
					codec = string(b[16:20]) // this should be the right codec
				}
			}
			panic(codec)
		}
	}
}
