package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	// get all context
	getCtxCmdOutput := new(bytes.Buffer)
	getCtxCmd := exec.Command("bash", "-c", "kubectl config get-contexts | awk 'NR!=1{print $2}' | xargs")
	getCtxCmd.Stdout = getCtxCmdOutput
	getCtxCmd.Stderr = os.Stderr
	err := getCtxCmd.Run()
	if err != nil {
		log.Fatalf("failed to run `kubectl config get-contexts | awk 'NR!=1{print $2}' | xargs` , err: %v", err)
	}
	if getCtxCmdOutput.String() == "\n" {
		log.Fatalf("no context detected")
	}

	// wait for selecting context
	prompt := promptui.Select{
		Label:        "switch context",
		Size:         10,
		Items:        strings.Split(string(strings.Replace(getCtxCmdOutput.String(), "\n", "", -1)), " "),
		HideSelected: true,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("failed to prepare selected menu, err: %v", err)
	}

	// switch context
	useCtxOutput := exec.Command("bash", "-c", fmt.Sprintf("kubectl config use-context %v", result))
	useCtxOutput.Stdout = os.Stdout
	useCtxOutput.Stderr = os.Stderr
	err = useCtxOutput.Run()
	if err != nil {
		log.Fatalf("failed to run `%v`, err: %v", fmt.Sprintf("kubectl config use-context %v", result), err)
	}
}
