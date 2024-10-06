package docs_blog_engine

import (
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
	// start the command
	if err := cmd.Start(); err != nil {
		return err
	}
	// read the standard output of the command
	fmt.Println("Building the app...")
	if err := cmd.Wait(); err != nil {
		return err
	}
	fmt.Println(stdout)
	fmt.Println("Build completed successfully.")
	return nil
}
