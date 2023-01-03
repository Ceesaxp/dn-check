# Domain Name Check

A simple tool to check if a domain name (or a list of names) is available under a list of TLDs. Can output results to either a console or a file.

## Usage

```bash
$ domain-check [options]
```

## Options:
- `-d string` Comma separated list of TLDs to check, defaults to com.
- `-f filename` File name to read the list of names from, one name per line. Superseded by the `-n` option.
- `-h`    Show help message and quit.
- `-j`    Output using JSON format.
- `-n string` List of names to check, separated by comma. Takes precedence over `-f` option.
- `-o filename` Spool output to a filename provided.
- `-v`    Enable verbose mode.

If `-d` is not specified, the script will check the `com` TLD. Either an `-f` or `-n` option must be specified. If both are specified, the `-n` option takes precedence.

## Example

```bash
$ domain-check -f domains.txt -o results.txt -d com,net,org -j

[will output the results to results.txt in JSON format]
```

```bash
$ domain-check -n yahoo,sun4everyone -d com,net,org,tj -v
Checking 2 names for 4 TLDs
Names        com  net  org  tj  
yahoo        NO   NO   NO   NO   
sun4everyone NO   YES  YES  YES  
```
