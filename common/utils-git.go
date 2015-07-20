package common

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func LocalGitBranch() string {
	gitRootDir, err := GitRootDir()
	if err != nil {
		return ""
	}

	b, err := exec.Command("git", "--git-dir", fmt.Sprintf("%s/.git", gitRootDir), "name-rev", "--name-only", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(string(b), "https://", "http://", -1))
}

func RemoteGitUrl() string {
	gitRootDir, err := GitRootDir()
	if err != nil {
		return ""
	}

	b, err := exec.Command("git", "--git-dir", fmt.Sprintf("%s/.git", gitRootDir), "config", "remote.origin.url").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(string(b), "https://", "http://", -1))
}

func GitRootDir() (string, error) {
	b, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	} else {
		return strings.TrimSpace(string(b)), err
	}
}

func PathRelativeToGitRoot(dirPath string) (string, error) {
	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return "", err
	}
	dir, err := os.Open(dirPath)
	if err != nil {
		return "", err
	}
	defer dir.Close()
	dirInfo, err := dir.Stat()
	if err != nil {
		return "", err
	}

	gitRootDir, err := GitRootDir()
	if err != nil {
		return "", err
	}
	root, err := os.Open(gitRootDir)
	if err != nil {
		return "", err
	}
	defer root.Close()
	rootInfo, err := root.Stat()
	if err != nil {
		return "", err
	}

	relativePath := ""
	for !os.SameFile(rootInfo, dirInfo) {
		relativePath = path.Base(dirPath) + "/" + relativePath
		dirPath = path.Dir(dirPath)
		dir, err = os.Open(dirPath)
		if err != nil {
			return "", err
		}
		defer dir.Close()
		dirInfo, err = dir.Stat()
		if err != nil {
			return "", err
		}
	}

	if relativePath == "" {
		relativePath = "."
	}

	return relativePath, nil
}
