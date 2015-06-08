package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func Main(do func() error) {
	if err := do(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func Cmd(stdout io.Writer, stderr io.Writer, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v: %v", args, err)
	}
	return nil
}

func ReadAll(filePath string) (retValue []byte, retErr error) {
	file, err := os.Open("Godeps/Godeps.json")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return ioutil.ReadAll(file)
}

func Println(object interface{}) {
	fmt.Printf("%+v\n", object)
}
