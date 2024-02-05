package utils

import (
	"io"
	"os"
	"os/exec"
)

func CallFFmpeg(stdin io.Reader, args ...string) error {
	args = append([]string{"-protocol_whitelist", "file,pipe,fd"}, args...)
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = cmd.Stdout
	err := cmd.Run()
	return err
}
