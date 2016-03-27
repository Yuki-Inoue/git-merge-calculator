
package main

import (
    "bufio"
    "fmt"
	"os/exec"
    "os"
    "io"
    "strings"
	"github.com/deckarep/golang-set"
)

func execCommand(progn string, args ...string) string {

	c, _ := exec.Command(progn, args...).Output()
	return string(c[:])
}

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s MergedCommit TargetBranch\n", os.Args[0])
		os.Exit(1)
	}

	mergedCommit := os.Args[1]
	targetBranch := os.Args[2]

	logRange := mergedCommit + ".." + targetBranch

	str := execCommand("git", "log", "--ancestry-path", "--pretty=format:%H", logRange)

	revs := strings.Split(str, "\n")
	set := mapset.NewSet()

	for _, v := range revs {
		set.Add(v)
	}

	// fmt.Printf(set.String())

	cmd := exec.Command("git", "rev-list", "--first-parent", targetBranch)
	out, _ := cmd.StdoutPipe()
	cmd.Start()

	bio := bufio.NewReader(out)

	answerByte, _ , _ := bio.ReadLine()
	answer := string(answerByte)

	for {

		line, hasMoreLine, err := bio.ReadLine()
		if !hasMoreLine && err == io.EOF {
			fmt.Println(err)
			break
		}

		rev := string(line[:])
		if !set.Contains(rev) {
			fmt.Println(answer)
			break
		} else {
			answer = rev
		}
	}

}
