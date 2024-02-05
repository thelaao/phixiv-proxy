package utils

import (
	"io"
	"os/exec"
)

func CallFFmpeg(stdin io.Reader, args ...string) error {
	args = append([]string{"-protocol_whitelist", "file,pipe,fd"}, args...)
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdin = stdin
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	return err
}
