package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

// authorErrorCount 用于排序作者错误统计
type authorErrorCount struct {
	Author string
	Count  int
}

// 更新后的错误解析正则表达式
// 匹配类似: path/to/file.ts(123,45): error TS2322: Message
var tsErrorRegex = regexp.MustCompile(`^([\w./\\~@%:-]+)\((\d+),(\d+)\):\serror\sTS\d+:`)

// git blame 输出中作者信息的正则表达式 (配合 --line-porcelain)
var gitAuthorRegex = regexp.MustCompile(`^author\s+(.*)`)

// ANSI 转义序列的正则表达式 (用于剥离颜色代码等)
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// stripAnsi 从字符串中移除 ANSI 转义序列
func stripAnsi(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

func main() {
	authorErrors := make(map[string]int)
	var processedErrorLines int
	var actualErrorsFoundByRegex int

	log.SetOutput(os.Stderr)
	log.Println("Starting ts-blame-reporter. Reading from stdin...")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rawLine := scanner.Text()
		log.Printf("SCRIPT_RECEIVED_RAW_LINE: %q", rawLine)

		line := stripAnsi(rawLine)
		if line != rawLine {
			log.Printf("SCRIPT_RECEIVED_STRIPPED_LINE (changed): %q", line)
		} else if line != "" {
			log.Printf("SCRIPT_RECEIVED_STRIPPED_LINE (no change): %q", line)
		}

		filePath, lineNumber, isError := parseErrorLine(line)

		if isError {
			actualErrorsFoundByRegex++
			log.Printf("REGEX_MATCHED: File: %q, Line: %q, FromStrippedLine: %q", filePath, lineNumber, line)

			author, err := getAuthorFromGitBlame(filePath, lineNumber)
			if err != nil {
				log.Printf("ERROR_GIT_BLAME: File: %q, Line: %q, Error: %v. Attributing to 'unknown_blame_error'.", filePath, lineNumber, err)
				author = "unknown_blame_error"
			}
			if author == "" {
				author = "unknown_author_empty"
			}
			authorErrors[author]++
		}
		processedErrorLines++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading stdin: %v", err)
	}

	log.Printf("Finished processing stdin. Total lines processed from input: %d. Lines matched by error regex: %d.", processedErrorLines, actualErrorsFoundByRegex)
	printReport(authorErrors)
}

// parseErrorLine 解析单行日志，提取文件路径和行号 (已更新)
func parseErrorLine(logLine string) (filePath string, lineNumber string, isError bool) {
	if logLine == "" {
		return "", "", false
	}
	matches := tsErrorRegex.FindStringSubmatch(logLine)
	// 新的正则表达式有3个捕获组 (文件路径, 行号, 列号)，所以完整匹配的len(matches)应该是4
	if len(matches) == 4 {
		filePath = strings.ReplaceAll(matches[1], `\`, `/`) // 捕获组1是文件路径
		lineNumber = matches[2]                             // 捕获组2是行号
		// matches[3] 是列号, 当前未使用
		return filePath, lineNumber, true
	}
	return "", "", false
}

// getAuthorFromGitBlame 执行 git blame 并解析作者信息
func getAuthorFromGitBlame(filePath string, lineNumber string) (author string, err error) {
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return "unknown_file_not_found", fmt.Errorf("file %s not found: %w", filePath, statErr)
	}

	cmd := exec.Command("git", "blame", "-L", fmt.Sprintf("%s,%s", lineNumber, lineNumber), "--line-porcelain", filePath)
	// log.Printf("DEBUG_GIT_CMD: %s", cmd.String())

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		errMsg := strings.ReplaceAll(strings.TrimSpace(stderr.String()), "\n", " ")
		return "git_blame_execution_error", fmt.Errorf("git blame failed for %s:%s: %w (stderr: %s)", filePath, lineNumber, err, errMsg)
	}

	blameScanner := bufio.NewScanner(strings.NewReader(out.String()))
	for blameScanner.Scan() {
		blameLine := blameScanner.Text()
		matches := gitAuthorRegex.FindStringSubmatch(blameLine)
		if len(matches) == 2 {
			return strings.TrimSpace(matches[1]), nil
		}
	}
	if err := blameScanner.Err(); err != nil {
		return "error_scanning_blame_output", fmt.Errorf("error scanning git blame output for %s:%s: %w", filePath, lineNumber, err)
	}

	return "unknown_author_not_found_in_blame", nil
}

// printReport 打印错误统计报告
func printReport(errorCounts map[string]int) {
	if len(errorCounts) == 0 {
		fmt.Println("\n--- TypeScript Error Report by Author ---")
		fmt.Println("No TypeScript errors matching the pattern were attributed to authors.")
		fmt.Println("---------------------------------------")
		log.Println("Report: No errors attributed to authors based on the provided input and regex.")
		return
	}

	counts := make([]authorErrorCount, 0, len(errorCounts))
	for author, count := range errorCounts {
		counts = append(counts, authorErrorCount{author, count})
	}

	sort.Slice(counts, func(i, j int) bool {
		if counts[i].Count == counts[j].Count {
			return counts[i].Author < counts[j].Author
		}
		return counts[i].Count > counts[j].Count
	})

	fmt.Println("\n--- TypeScript Error Report by Author ---")
	totalReportedErrors := 0
	for _, item := range counts {
		fmt.Printf("%-40s: %d errors\n", item.Author, item.Count)
		totalReportedErrors += item.Count
	}
	fmt.Println("---------------------------------------")
	fmt.Printf("Total TypeScript errors attributed: %d\n", totalReportedErrors)
	log.Printf("Report: Generated for %d authors with a total of %d errors attributed.", len(counts), totalReportedErrors)
}
