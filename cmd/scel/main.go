package main

import (
	"io/ioutil"
	"log"

	"encoding/json"

	"fmt"

	"flag"

	"bitbucket.org/johncming/scel"
)

var scelPath string

func init() {
	flag.StringVar(&scelPath, "p", "", "scel path")
	flag.Parse()
}

func readScel() ([]byte, error) {
	return ioutil.ReadFile(scelPath)
}

func main() {
	data, err := readScel()
	if err != nil {
		log.Fatalln(err)
	}

	scel := scel.NewScel(data)
	err = scel.Run()
	if err != nil {
		log.Fatalln(err)
	}

	output, err := json.Marshal(scel.WordPy)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s\n", string(output))
}
