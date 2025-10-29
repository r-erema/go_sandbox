//nolint:gosec // disable G204
package utils

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"path"
)

func GoCompileAndStart(
	ctx context.Context,
	file string,
	args ...string,
) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	binFile := path.Join("/", "tmp", Hash(file))

	if err := goCompile(file, binFile); err != nil {
		return nil, nil, nil, fmt.Errorf("go compile failed: %w", err)
	}

	cmd := exec.CommandContext(ctx, binFile, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("exec: StdoutPipe failed: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("exec: StderrPipe failed: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, nil, fmt.Errorf("start %s failed: %w", binFile, err)
	}

	return cmd, stdout, stderr, nil
}

func GoCompileAndStartInNetNs(
	ctx context.Context,
	file, netNs string,
	args ...string,
) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	binFile := path.Join("/", "tmp", Hash(file))

	if err := goCompile(file, binFile); err != nil {
		return nil, nil, nil, fmt.Errorf("go compile failed: %w", err)
	}

	cmdArgs := append([]string{"netns", "exec", netNs, binFile}, args...)
	cmd := exec.CommandContext(ctx, "ip", cmdArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("exec: StdoutPipe failed: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("exec: StderrPipe failed: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, nil, fmt.Errorf("exec: failed to start cmd: %w", err)
	}

	return cmd, stdout, stderr, nil
}

func goCompile(file, binFile string) error {
	if err := exec.Command("go", "build", "-o", binFile, file).Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	return nil
}
