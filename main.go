package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"strings"
)

const usage = `Usage: dn-check [options] [file ...]
Options:
  -f file : The file containing the domain names to check.
  -o file : The file to write the results to.
  -d list : A comma-separated list of TLDs to check
  -h : Show the help message.
`

// Options struct
type Options struct {
	FileName string // File name to read the list of names from
	TLDs     string // List of TLDs to check
	Output   string // Output file name
	Verbose  bool   // Verbose mode
	Help     bool   // Help

	// Extracted from CLI options
	TLDsList  []string
	NamesList []string
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
func readNamesFromFile(filename string) ([]string, error) {
	d, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	} else {
		return strings.Split(strings.ToLower(string(d)), "\n"), nil
	}
}

func readOptions() Options {
	var options Options
	flag.StringVar(&options.FileName, "f", "", "File name to read the list of names from, one name per line.")
	flag.StringVar(&options.TLDs, "d", "com", "Comma separated list of TLDs to check")
	flag.StringVar(&options.Output, "o", "", "Spool output to a `filename` provided")
	flag.BoolVar(&options.Verbose, "v", false, "Enable `verbose` mode") // not done
	flag.BoolVar(&options.Help, "h", false, "Help message")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	// Before anything else – check if we're asked for help
	if options.Help || options.FileName == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// put TLDs into a list from the comma separated string
	options.TLDsList = strings.Split(options.TLDs, ",")
	// read domain names from file specified via CLI and add them to the list
	l, err := readNamesFromFile(options.FileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		options.NamesList = l
	}

	return options
}

func run(opts Options) {
	for _, name := range opts.NamesList {
		if name == "" {
			// skip empty lines
			continue
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
	}
}

func main() {
	opts := readOptions()
	run(opts)
}
