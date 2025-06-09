
# ts-blame-reporter

`ts-blame-reporter` is a command-line tool that analyzes TypeScript (or similarly formatted) compiler error logs. It uses `git blame` to trace each error to the last person who modified the related line of code, then summarizes this data in a report.

This tool is intended to improve error ownership, facilitate code reviews, and raise awareness of code quality across the team.

---

## Features

- Reads compiler output from standard input (stdin)
- Parses error lines using a regular expression matching `filepath(line,column): error TSxxxx: message`
- Calls `git blame` for each matched line to identify the responsible author
- Outputs a summary of errors per author in descending order
- Strips ANSI escape sequences (like color codes) from logs to avoid parse issues

---

## Prerequisites

- **Go**: Version 1.18 or higher (for compiling `ts-blame-reporter`)
- **Git**: Must be installed and available in the system's `PATH`
- **Node.js project**: You should run the tool in a project using TypeScript

---

## Installation

### 1. Get the Code

```bash
git clone https://github.com/tomwu618/ts-blame-reporter.git
cd ts-blame-reporter
```

Alternatively, if you only have the `main.go` file, ensure it is saved in your current working directory.

### 2. Compile

```bash
go build -o ts-blame-reporter main.go
```

This will generate an executable named `ts-blame-reporter` (or `ts-blame-reporter.exe` on Windows).

#### Optional: Cross-compiling for another platform

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o ts-blame-reporter main.go

# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o ts-blame-reporter.exe main.go

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o ts-blame-reporter main.go
```

---

## Usage

Before using this tool, make sure you’ve installed the necessary dependencies for your project:

```bash
npm install
npm install -D vue-tsc
```

### macOS / Linux

Move the binary to a directory included in your system's `PATH`:

```bash
sudo mv ts-blame-reporter /usr/local/bin/
sudo chmod +x /usr/local/bin/ts-blame-reporter
```

Run your TypeScript build command and pipe its output to `ts-blame-reporter`:

```bash
npm run build:dev 2>&1 | ts-blame-reporter
```

To redirect the report to a file:

```bash
npm run build:dev 2>&1 | ts-blame-reporter > error_report.txt
```

---

### Windows

Place the compiled binary (`ts-blame-reporter.exe`) into a directory like `C:\Windows\System32` (included in your `PATH`), or use the full path each time.

Then run:

```powershell
npm run build:dev 2>&1 | Out-String | C:\Windows\System32\ts-blame-reporter.exe
```

If needed, add `.exe` to `PATH` or create a custom alias for convenience.

---

## Example Output

```
--- TypeScript Error Report by Author ---
songjunhao                              : 101 errors
wangzhengxue                            : 93 errors
chinux_2012                             : 8 errors
yangmiaomiao                            : 2 errors
unknown_blame_error                     : 1 errors
---------------------------------------
Total TypeScript errors attributed: 205
```

---

## How It Works

1. **Read Input**: Reads the compiler/build logs line-by-line from stdin
2. **Strip ANSI Codes**: Cleans up log lines by removing terminal color codes
3. **Parse Error Lines**: Matches lines in the format `file(line,column): error TSxxxx: message`
4. **Blame Lookup**: Runs `git blame` to determine who last changed the error line
5. **Aggregate Report**: Summarizes how many errors are attributed to each author

---

## Known Issues / Limitations

- **Blame Accuracy**: `git blame` attributes changes to the last modifier, not necessarily the original author
- **File Path Mismatches**: Log paths must match repository paths exactly (case, separators, etc.)
- **Performance**: Many errors may slow the tool down due to many individual `git blame` calls
- **Regex Limitation**: Only supports one hardcoded error format for now (TypeScript)

---

## Contributing

Pull Requests and Issues are welcome!

---

## License

MIT License. See `LICENSE` file for details.
