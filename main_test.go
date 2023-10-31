package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"testing"
)

func checkExpectation(t *testing.T, testDir string) {
	want, err := os.ReadFile(path.Join(testDir, "want.yaml"))
	if err != nil {
		t.Fatalf("error reading expected resources file %s: %v", path.Join(testDir, "want.yaml"), err)
	}

	arg := []string{"build"}
	flags := []string{"--enable-alpha-plugins", "--enable-exec"}
	arg = append(arg, flags...)
	arg = append(arg, testDir)

	cmd := exec.Command("kustomize", arg...)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()

	got := out.Bytes()
	if bytes.Compare(want, got) != 0 {
		t.Errorf("wanted %s. got %s", want, got)
	}
}

func TestKasque_KRM(t *testing.T) {
}
