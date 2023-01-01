package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"strings"
)

// Options struct
type Options struct {
	// File name to read the list of names from
	FileName string
	// List of TLDs to check
	TLDs string
	// Output file name
	Output string
	// Verbose mode
	Verbose bool
	// Extarcted from CLI options
	TLDsList  []string
	NamesList []string
	// Help
	Help bool
}

func isDomainNameAvailable(domain string) (bool, error) {
	_, err := net.LookupHost(domain)
	if err != nil {
		if _, ok := err.(*net.DNSError); ok {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// Open filename for read and read all lines into a list
// Returns the list of names
func readNamesFromFile(filename string) []string {
	d, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l := strings.Split(string(d), "\n")
	if l[len(l)-1] == "" {
		// drop last element if it is blank
		return l[0 : len(l)-1]
	} else {
		return l
	}
}

// Read options from the command line:
// -f filename - read the list of names from the file
// -d tlds - comma separated list of TLDs to check
// -n names - comma separated list of names to check
// -o output - output file name
// -v - verbose mode
// -h - help
func readOptions() Options {
	var options Options
	flag.StringVar(&options.FileName, "f", "", "File name to read the list of names from, one name per line.")
	flag.StringVar(&options.TLDs, "d", "com,ru,tech,ai", "Comma separated list of TLDs to check")
	flag.StringVar(&options.Output, "o", "", "Spool output to a `filename` provided")
	flag.BoolVar(&options.Verbose, "v", false, "Enable `verbose` mode") // not done
	flag.BoolVar(&options.Help, "h", false, "Help message")
	flag.Parse()

	// Before anything else – check if we're asked for help
	if options.Help || options.FileName == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// put TLDs into a list from the comma separated string
	options.TLDsList = strings.Split(options.TLDs, ",")
	// read domain names from file specified via CLI and add them to the list
	options.NamesList = readNamesFromFile(options.FileName)

	return options
}

func main() {
	opts := readOptions()

	for _, name := range opts.NamesList {
		if name == "" {
			goto LAST
		}
		for _, tld := range opts.TLDsList {
			dn, err := isDomainNameAvailable(name + "." + tld)
			if err != nil {
				fmt.Println(err)
			}
			if dn {
				fmt.Printf("%s.%s is available\n", name, tld)
			} else {
				fmt.Printf("%s.%s is %s available\n", name, tld, color.RedString("NOT"))
			}
		}
	LAST:
	}

	fmt.Printf("OK")
}
