package virtualbox

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

var (
	VBM     = "VBoxManage" // Path to VBoxManage utility.
	Verbose = new(bool)    // Verbose mode.
)

func init() {
	if p := os.Getenv("VBOX_INSTALL_PATH"); p != "" && runtime.GOOS == "windows" {
		VBM = filepath.Join(p, "VBoxManage.exe")
	}
}

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
	ErrVBMNotFound     = errors.New("VBoxManage not found")
)

func vbm(args ...string) error {
	cmd := exec.Command(VBM, args...)
	if *Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Printf("exec %s with arguments %v", VBM, args)
	}
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			return ErrVBMNotFound
		}
		return err
	}
	return nil
}

func vbmOut(args ...string) ([]byte, error) {
	cmd := exec.Command(VBM, args...)
	if *Verbose {
		cmd.Stderr = os.Stderr
		log.Printf("exec %s with arguments %v", VBM, args)
	}

	b, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return b, err
}

func vbmOutErr(args ...string) ([]byte, []byte, error) {
	cmd := exec.Command(VBM, args...)
	if *Verbose {
		log.Printf("exec %s with arguments %v", VBM, args)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return stdout.Bytes(), stderr.Bytes(), err
}