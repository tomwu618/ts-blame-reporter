# ts-blame-reporter

`ts-blame-reporter` is a command-line tool designed to analyze TypeScript (or other similarly formatted) build error logs. It integrates with `git blame` to attribute each detected error to the developer who last modified the relevant line of code, generating a simple statistical report.

This tool aims to help development teams trace the origin of errors, assist in code reviews, and enhance overall awareness of code quality.

## Features

* Reads build logs from standard input (stdin).
* Parses error lines using a regular expression (currently hardcoded to match `filepath(line,column): error TSxxxx: message` format).
* Invokes `git blame` for each parsed error line to retrieve the author who last modified that line.
* Counts and reports the number of errors attributed to each author, sorted descending by error count.
* Attempts to strip ANSI color codes from log lines to improve parsing accuracy.

## Prerequisites

* **Go**: Version 1.18 or later (for compiling `ts-blame-reporter`).
* **Git**: Must be installed and available in the system's `PATH`. `ts-blame-reporter` executes `git` commands within the project directory being analyzed.

## Installation

1.  **Get the Code**:
    Clone this repository:
    ```bash
    git clone [https://github.com/tomwu618/ts-blame-reporter.git](https://github.com/tomwu618/ts-blame-reporter.git)
    cd ts-blame-reporter
    ```
    Alternatively, if you only have the `main.go` file, ensure it's in your current working directory.

2.  **Compile**:
    In the project root directory (where `main.go` is located), run the following command to compile:
    ```bash
    go build -o ts-blame-reporter main.go
    ```
    This will generate an executable file named `ts-blame-reporter` (or `ts-blame-reporter.exe` on Windows).

    * **Cross-compiling for a specific platform (e.g., macOS Intel)**:
        ```bash
        # macOS Intel (amd64)
        GOOS=darwin GOARCH=amd64 go build -o ts-blame-reporter main.go
        # Windows (amd64)
        GOOS=windows GOARCH=amd64 go build -o ts-blame-reporter.exe main.go
        # Linux (amd64)
        GOOS=linux GOARCH=amd64 go build -o ts-blame-reporter main.go
        ```

3.  **Place in System PATH**:
    Move the compiled executable `ts-blame-reporter` to a directory included in your system's `PATH` environment variable. This allows you to call it directly from any path.
    * On macOS or Linux, a common directory is `/usr/local/bin`:
        ```bash
        sudo mv ts-blame-reporter /usr/local/bin/
        sudo chmod +x /usr/local/bin/ts-blame-reporter # Ensure execute permission
        ```
    * On Windows, you can place it in a directory that is already part of your `Path` environment variable.

## Usage

In the root directory of your TypeScript (or other) project, pipe the output of your build/compilation command to `ts-blame-reporter`. **Ensure you also redirect the standard error stream (`stderr`), as many build tools output errors to `stderr`.**

```bash
<your-build-command> 2>&1 | ts-blame-reporter


Example (for a Vue + TypeScript project):

Bash


npm run build:dev 2>&1 | ts-blame-reporter


Alternatively, if your build command only outputs errors to stdout:

Bash


npm run build:dev | ts-blame-reporter


Output
The tool outputs debug information to standard error (stderr) and the final statistical report to standard output (stdout).
You can redirect the report to a file:

Bash


npm run build:dev 2>&1 | ts-blame-reporter > error_report.txt


Example Output Report:



--- TypeScript Error Report by Author ---
songjunhao                              : 101 errors
wangzhengxue                            : 93 errors
chinux_2012                             : 8 errors
yangmiaomiao                            : 2 errors
unknown_blame_error                     : 1 errors
---------------------------------------
Total TypeScript errors attributed: 205


How it Works
Read Input: Reads logs line by line from standard input.
Strip ANSI Codes: Removes potential ANSI SGR color codes from each line.
Parse Error Lines: Uses a built-in regular expression to match lines fitting the filepath(line,column): error TSxxxx: message pattern, extracting the file path and line number.
Get Author: For each matched error, it executes the git blame -L <line_number>,<line_number> --line-porcelain <file_path> command.
Parse Blame Output: Extracts the author field from the git blame output.
Aggregate & Report: Counts the number of errors attributed to each author and prints the final report, sorted descending by error count.
Known Issues / Limitations
git blame Attribution: git blame points to the "last modifier" of a line of code. This may not always be the original author who introduced the error (e.g., a code formatting commit might alter many lines and be incorrectly attributed). Use the results of this tool as an aid for analysis.
File Path Matching: If the file paths reported in the build log do not exactly match the actual paths tracked in the Git repository (including case sensitivity or relative/absolute path differences), git blame might fail (resulting in unknown_file_not_found or git_blame_execution_error in the logs).
Performance: For very large logs with a vast number of errors (thousands+), performance might degrade due to sequential git blame calls for each error. Future optimizations could include batching blame calls per file or concurrent execution.
Hardcoded Regex: The regular expression for parsing error lines is currently hardcoded. If your project's error format differs, you'll need to modify the tsErrorRegex variable in the source code.
Contributing
Pull Requests and Issues are welcome to improve this tool!
License
This project is licensed under the MIT License - see the LICENSE file for details.
