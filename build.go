package docs_blog_engine

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Build() error {
	// define the command that you want to run
	cmd := exec.Command("npm", "run", "build")
	// specify the working directory of the command
	cmd.Dir = "./app/"
	// get the pipe for the standard output of the command
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// get the pipe for the standard error of the command
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create a scanner to read from stdout
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // print each line from stdout
		}
	}()

	// Create a scanner to read from stderr (for error messages)
	errScanner := bufio.NewScanner(stderr)
	go func() {
		for errScanner.Scan() {
			fmt.Println(errScanner.Text()) // print each line from stderr
		}
	}()

	// wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Println("Build completed successfully.")
	return nil
}

func Install() error {
	// define the command that you want to run
	cmd := exec.Command("npm", "run", "install")
	// specify the working directory of the command
	cmd.Dir = "./app/"
	// get the pipe for the standard output of the command
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// get the pipe for the standard error of the command
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create a scanner to read from stdout
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // print each line from stdout
		}
	}()

	// Create a scanner to read from stderr (for error messages)
	errScanner := bufio.NewScanner(stderr)
	go func() {
		for errScanner.Scan() {
			fmt.Println(errScanner.Text()) // print each line from stderr
		}
	}()

	// wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Println("Build completed successfully.")
	return nil
}
