package docs_blog_engine_run

import (
	"bufio"
	"fmt"
	"os/exec"
)

func runCommand(path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = path
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // print each line from stdout
		}
	}()

	errScanner := bufio.NewScanner(stderr)
	go func() {
		for errScanner.Scan() {
			fmt.Println(errScanner.Text()) // print each line from stderr
		}
	}()

	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Println("Command ran successfully.")
	return nil
}
