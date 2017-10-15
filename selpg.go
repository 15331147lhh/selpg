package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	flagSet = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	start   = flagSet.Int("s", 0, "The start page number")
	end     = flagSet.Int("e", 0, "The end page number")
	l       = flagSet.Int("l", 72, "The number of line in one page")
	f       = flagSet.Bool("f", false, "read one page until '\f' ")
)

func printError(message string) {
	err := errors.New(message)
	fmt.Fprintln(os.Stderr, "error>>", err)
}

func checkForSE() bool {
	if len(os.Args[0]) <= 2 {
		printError("you need exc in this form : './selpg  -s=?  -e=?  filename'")
		printError("-s and -e option are both in need")
		return false
	}
	if os.Args[1][0:2] != "-s" {
		printError("-s should be first option")
		return false
	}
	if os.Args[2][0:2] != "-e" {
		printError("-e should be second option")
		return false
	}
	if *start <= 0 || *end <= 0 {
		printError("-s and -e should be bigger than 0")
		return false
	}
	if *start > *end {
		printError("-s should be smaller than -e")
		return false
	}
	if *l <= 0 {
		printError("-l should be bigger than 0")
		return false
	}
	return true
}

func fileIO(Ibuf *bufio.Reader, Obuf *os.File) {
	count := *end - *start + 1
	if !*f {
		for i := 1; i < *start; i++ {
			for j := 0; j < *l; j++ {
				Ibuf.ReadString('\n')
			}
		}
		for i := 0; i < count; i++ {
			for j := 0; j < *l; j++ {
				line, err := Ibuf.ReadString('\n')
				if err != nil {
					if err == io.EOF && i != count && j != *l {
						printError("no enough page of the file")
						return
					} else {
						fmt.Fprint(os.Stderr, "error>>", err.Error())
					}
				}
				if Obuf != nil {
					Obuf.WriteString(line)
				} else {
					fmt.Print(line)
				}
			}
		}
	} else {
		for i := 1; i < *start; i++ {
			Ibuf.ReadString('\f')
		}
		for i := 0; i < count; i++ {
			line, err := Ibuf.ReadString('\f')
			if err != nil {
				if err == io.EOF && i != count {
					printError("no enough page of the file")
					return
				} else {
					fmt.Fprint(os.Stderr, "error>>", err.Error())
				}
			}
			if Obuf != nil {
				Obuf.WriteString(line)
			} else {
				fmt.Print(line)
			}
		}
	}
}

func iodata(inputFileName string, outputFileName string) {
	var Ibuf *bufio.Reader
	if inputFileName != "" {
		inputFile, err := os.OpenFile(inputFileName, os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error>>", err.Error())
		}
		Ibuf = bufio.NewReader(inputFile)
	} else {
		Ibuf = bufio.NewReader(os.Stdin)
	}
	var Obuf *os.File
	var err error
	if outputFileName != "" {
		Obuf, err = os.OpenFile(outputFileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error>>", err.Error())
		}
	} else {
		Obuf = nil
	}
	fileIO(Ibuf, Obuf)
}

func main() {
	flagSet.Parse(os.Args[1:])
	if checkForSE() {
		var inputFile = ""
		var outputFile = ""
		if flagSet.NArg() > 0 {
			inputFile = flagSet.Arg(0)
		}
		if flagSet.NArg() > 1 {
			outputFile = flagSet.Arg(1)
		}
		iodata(inputFile, outputFile)
	}

}
