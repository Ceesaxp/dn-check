package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"strings"
)

const usage = `Usage: dn-check [options]
Options:
  -d string
        Comma separated list of TLDs to check, defaults to com.
  -f filename
        File name to read the list of names from, one name per line. Superseded by the -n option.
  -h    Show this help message.
  -j    Output using JSON format.
  -n string
        List of names to check, separated by comma. Takes precedence over -f option.
  -o filename
        Spool output to a filename provided.
  -v    Enable verbose mode.
`

// Options structure to hold the command line options
type Options struct {
	FileName  string   // FileName to read the list of names from ... OR
	Names     string   // List of Names to check, separated by comma (takes precedence over FileName)
	TLDs      string   // List of TLDs to check
	Output    string   // Output file name
	Verbose   bool     // Verbose mode
	Json      bool     // Output JSON
	Help      bool     // Help
	TLDsList  []string // List of TLDs to check
	NamesList []string // List of names either from the command line or read from file
}

// TLD structure to keep the top-level domain availability data. Will be nested in the Result structure
type TLD struct {
	TLDName     string `json:"top_level_domain"`
	IsAvailable bool   `json:"is_available"`
}

// Result structure to keep the results of the check
type Result struct {
	Name    string `json:"domain_name"`
	TLDList []TLD  `json:"tlds"`
}

// Lookup a domain name using. Return true if the name is available, otherwise false
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

// Read the command line options
func readOptions() Options {
	var options Options
	flag.StringVar(&options.FileName, "f", "", "File name to read the list of names from, one name per line. Superseded by the -n option.")
	flag.StringVar(&options.Names, "n", "", "List of names to check, separated by comma. Takes precedence over -f option.")
	flag.StringVar(&options.TLDs, "d", "com", "Comma separated list of TLDs to check.")
	flag.StringVar(&options.Output, "o", "", "Spool output to a `filename` provided.")
	flag.BoolVar(&options.Json, "j", false, "Output using JSON format.")
	flag.BoolVar(&options.Verbose, "v", false, "Enable verbose mode.")
	flag.BoolVar(&options.Help, "h", false, "Show this help message.")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	// Before anything else ??? check if we're asked for help
	if options.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// put TLDs into a list from the comma separated string
	options.TLDsList = strings.Split(options.TLDs, ",")

	if options.Names == "" {
		if options.FileName == "" {
			fmt.Println("Error: No names provided")
			flag.PrintDefaults()
			os.Exit(1)
		} else {
			// Read names from file
			options.NamesList, _ = readNamesFromFile(options.FileName)
		}
	} else {
		// Read names from command line
		options.NamesList = strings.Split(options.Names, ",")
	}

	return options
}

func run(opts Options) ([]Result, error) {
	var Results []Result
	var tlds []TLD

	if opts.Verbose {
		PrintVerboseHeader(opts)
	}

	for _, name := range opts.NamesList {
		if name == "" {
			// skip empty lines
			continue
		}
		if opts.Verbose {
			fmt.Printf("%-12s ", name)
		}
		tlds = nil
		for _, tld := range opts.TLDsList {
			dnAvailable, err := isDomainNameAvailable(name + "." + tld)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			tlds = append(tlds, TLD{tld, dnAvailable})
			if opts.Verbose {
				VerboseOutput(dnAvailable)
			}
		}
		if opts.Verbose {
			fmt.Println()
		}
		Results = append(Results, Result{Name: name, TLDList: tlds})
	}
	return Results, nil
}

// VerboseOutput : Helper function to print YES/NO
func VerboseOutput(s bool) {
	if s {
		fmt.Print(color.GreenString("YES  "))
	} else {
		fmt.Print(color.RedString("NO   "))
	}
}

// PrintVerboseHeader : Prints the header for the verbose output
func PrintVerboseHeader(opts Options) {
	fmt.Println("Checking", len(opts.NamesList), "names for", len(opts.TLDsList), "TLDs")
	fmt.Print("Names       ")
	for _, t := range opts.TLDsList {
		fmt.Printf(" %-4s", t)
	}
	fmt.Println()
}

// SpoolOutputToFile : Output results to a file
func SpoolOutputToFile(outputFileName string, results []Result, jsonOutput bool) error {
	f, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}(f)
	if jsonOutput {
		// Format and output JSON into file if requested
		jsonText, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}
		_, err = f.Write(jsonText)
		if err != nil {
			return err
		}
	} else {
		for _, result := range results {
			for _, tld := range result.TLDList {
				_, err = fmt.Fprintf(f, "%s.%s : %t\n", result.Name, tld.TLDName, tld.IsAvailable)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	opts := readOptions()
	NamesResult, err := run(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		if opts.Output != "" {
			err := SpoolOutputToFile(opts.Output, NamesResult, opts.Json)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
