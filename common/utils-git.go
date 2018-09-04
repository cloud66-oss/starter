package common

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func HasGit(dir string) bool {
	_, err := GitRootDir(dir)
	if err != nil {
		return false
	}
	return true
}

func LocalGitBranch(dir string) string {
	gitRootDir, err := GitRootDir(dir)
	if err != nil {
		return ""
	}

	b, err := exec.Command("git", "--git-dir", fmt.Sprintf("%s/.git", gitRootDir), "name-rev", "--name-only", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(string(b), "https://", "http://", -1))
}

func AddFile(dir string, file string) error {
	gitRootDir, err := GitRootDir(dir)
	if err != nil {
		return err
	}

	_, err = exec.Command("git", "--git-dir", fmt.Sprintf("%s/.git", gitRootDir), "add", file).Output()
	if err != nil {
		return err
	}

	return nil
}

func Commit(dir string, message string) error {
	gitRootDir, err := GitRootDir(dir)
	if err != nil {
		return err
	}

	_, err = exec.Command("git", "--git-dir", fmt.Sprintf("%s/.git", gitRootDir), "commit", "-m", fmt.Sprintf("'%s'", message)).Output()
	if err != nil {
		return err
	}

	return nil
}

func RemoteGitUrl(dir string) string {
	gitRootDir, err := GitRootDir(dir)
	if err != nil {
		return ""
	}

	b, err := exec.Command("git", "--git-dir", fmt.Sprintf("%s/.git", gitRootDir), "config", "remote.origin.url").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(string(b), "https://", "http://", -1))
}

func GitRootDir(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	b, err := cmd.Output()
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

	gitRootDir, err := GitRootDir(dirPath)
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
