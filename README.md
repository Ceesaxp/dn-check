# Domain Name Check

A simple tool to check if a domain name is available.

## Usage

```bash
$ ./domain-check.sh [options]
```

## Options

- `-f` : The file containing the domain names to check.
- `-o` : The file to write the results to.
- `-h` : Show the help message.
- `-d` : A comma-separated list of TLDs to check

If `-d` is not specified, the script will check the `.com` TLD.

## Example

```bash
$ ./domain-check.sh -f domains.txt -o results.txt -d com,net,org
[will output the results to results.txt]
```