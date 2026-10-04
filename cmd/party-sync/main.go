// Command party-sync 跑队伍状态版本向量样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"example.com/partysync"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("party-sync", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "idem", "idem / gap / diff / missing / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *sample == "idem" {
		s := partysync.New()
		s.Apply("a", 1, map[string]string{"v": "v1"})
		second := s.Apply("a", 1, map[string]string{"v": "v2"})
		fmt.Fprintf(stdout, "state=%s second=%v\n", s.State("a")["v"], second)
		return 0
	}
	if *sample == "gap" {
		s := partysync.New()
		s.Apply("a", 1, map[string]string{"v": "m1"})
		s.Apply("a", 3, map[string]string{"v": "m3"})
		s.Apply("a", 2, map[string]string{"v": "m2"})
		fmt.Fprintf(stdout, "state=%s version=%d pending=%d\n",
			s.State("a")["v"], s.Version()["a"], s.Pending("a"))
		return 0
	}
	if *sample == "diff" {
		s := partysync.New()
		s.Apply("a", 1, map[string]string{"v": "m1"})
		s.Apply("b", 5, map[string]string{"v": "m5"})
		fmt.Fprintf(stdout, "diff=%s\n", strings.Join(s.Diff(partysync.Vector{"a": 1, "b": 3}), ","))
		return 0
	}
	if *sample == "missing" {
		s := partysync.New()
		s.Apply("a", 1, map[string]string{"v": "m1"})
		fmt.Fprintf(stdout, "missing=%s\n", strings.Join(s.Missing(partysync.Vector{"a": 1, "b": 4}), ","))
		return 0
	}
	if *sample == "scan" {
		s := partysync.New()
		for index := 0; index < 3000; index++ {
			s.Apply(fmt.Sprintf("m-%d", index), 1, map[string]string{"v": "x"})
		}
		fmt.Fprintf(stdout, "scanned=%d\n", s.Scanned())
		return 0
	}
	fmt.Fprintln(stderr, "需要 --sample idem|gap|diff|missing|scan")
	return 2
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
