package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

func runTerraformConsole(expression string) (string, error) {
	cmd := exec.Command("terraform", "console")
	cmd.Stdin = strings.NewReader(expression)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func runBat(inputText string) error {
	cmd := exec.Command("bat", "-fplterraform", "--pager", "less -R -+SG --mouse", "--theme", "Dracula")
	cmd.Stdin = strings.NewReader(inputText)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeExpression(expression string) {
	output, err := runTerraformConsole(expression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
	err = runBat(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running bat command: %s", err)
	}
}

func main() {
	if len(os.Args) > 1 {
		// Join command line arguments into a single expression
		expression := strings.Join(os.Args[1:], " ")
		executeExpression(expression)
	} else {
		// REPL mode
		rl, err := readline.New("> ")
		if err != nil {
			panic(err)
		}
		defer rl.Close()

		fmt.Println("Enter expressions. Type 'exit' or Ctrl-D to exit the REPL.")
		for {
			line, err := rl.Readline()
			if err == readline.ErrInterrupt || err == io.EOF {
				break
			}
			if strings.TrimSpace(line) == "exit" {
				break
			}
			executeExpression(line)
		}
	}
}
