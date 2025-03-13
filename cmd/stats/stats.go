package stats

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	StartCommand = &cobra.Command{
		Use:   "stats",
		Short: "code_stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	path string
)

const (
	fileTypeSource = "source"
	fileTypeTest   = "test"
)

func init() {
	StartCommand.Flags().StringVarP(&path, "path", "p", ".", "the source code dir")
}

func run() error {
	csTotal := &CodeStats{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		cs, err := countGolang(path)
		if err != nil {
			return err
		}
		csTotal.Add(cs)
		return nil
	})
	csTotal.Print()
	return err
}

func countGolang(path string) (CodeStats, error) {

	var (
		err error
		cs  CodeStats
		tp  = fileTypeSource
	)

	if filepath.Ext(path) != ".go" {
		return cs, err
	}
	if strings.HasSuffix(path, "_test.go") {
		cs.TestFileTotal++
		tp = fileTypeTest
	} else {
		cs.CodeFileTotal++
	}

	fl, err := os.Open(path)
	if err != nil {
		return cs, err
	}

	sc := bufio.NewScanner(fl)
	for sc.Scan() {
		text := sc.Text()
		if text == "\n" {
			//空行
			continue
		}

		switch tp {
		case fileTypeSource:
			if strings.HasPrefix(strings.TrimPrefix(text, " "), "//") {
				cs.CodeCommentTotal++
				continue
			}
			cs.CodeLineTotal++
		case fileTypeTest:
			cs.TestLineTotal++

		}
	}
	return cs, err
}

type CodeStats struct {
	CodeFileTotal    int `json:"code_file_total,omitempty"` //源码文件数据
	CodeLineTotal    int `json:"code_line_total,omitempty"` //源码行数，去除掉空行和注释行
	CodeCommentTotal int `json:"comment_total,omitempty"`   //注释行数
	TestFileTotal    int `json:"test_file_total,omitempty"` //测试文件数量
	TestLineTotal    int `json:"test_line_total,omitempty"` //测试代码行数，去掉空行和注释行
}

func (c *CodeStats) Add(target ...CodeStats) {
	for _, v := range target {
		c.CodeFileTotal += v.CodeFileTotal
		c.CodeLineTotal += v.CodeLineTotal
		c.CodeCommentTotal += v.CodeCommentTotal
		c.TestFileTotal += v.TestFileTotal
		c.TestLineTotal += v.TestLineTotal
	}
}

func (c *CodeStats) Print() {
	os.Stdout.WriteString(fmt.Sprintf("CODE_FILE_TOTAL=%d\n", c.CodeFileTotal))
	os.Stdout.WriteString(fmt.Sprintf("CODE_LINE_TOTAL=%d\n", c.CodeLineTotal))
	os.Stdout.WriteString(fmt.Sprintf("CODE_COMMENT_TOTAL=%d\n", c.CodeCommentTotal))
	os.Stdout.WriteString(fmt.Sprintf("TEST_FILE_TOTAL=%d\n", c.TestFileTotal))
	os.Stdout.WriteString(fmt.Sprintf("TEST_LINE_TOTAL=%d\n", c.TestLineTotal))
}

func closeStdout() {
	f := os.NewFile(0, os.DevNull)
	os.Stderr = f
}
func resetStdout() {
	os.Stdout = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	os.Stderr = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
}
