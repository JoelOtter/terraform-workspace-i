package terraform

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Workspace struct {
	Name    string
	Current bool
}

func GetWorkspaces() ([]Workspace, error) {
	tfCmd := exec.Command("terraform", "workspace", "list")
	tfCmd.Stderr = os.Stderr
	output, err := tfCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get terraform workspaces: %w", err)
	}
	var workspaces []Workspace
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		workspaceText := scanner.Text()
		workspaces = append(workspaces, Workspace{
			Name:    strings.TrimSpace(strings.TrimPrefix(workspaceText, "*")),
			Current: strings.HasPrefix(workspaceText, "*"),
		})
	}
	return workspaces, nil
}

func ChangeWorkspace(workspace string, output io.Writer) error {
	tfCmd := exec.Command("terraform", "workspace", "select", workspace)
	tfCmd.Stdout = output
	tfCmd.Stderr = output
	if err := tfCmd.Run(); err != nil {
		return fmt.Errorf("failed to select workspace %s: %w", workspace, err)
	}
	return nil
}

func DeleteWorkspace(workspace string, output io.Writer) ([]Workspace, error) {
	tfCmd := exec.Command("terraform", "workspace", "delete", workspace)
	tfCmd.Stdout = output
	tfCmd.Stderr = output
	if err := tfCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to delete workspace %s: %w", workspace, err)
	}
	return GetWorkspaces()
}
