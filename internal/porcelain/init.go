package porcelain

import (
	"os"
	"path/filepath"
)

func Init(path string) error {
	gotDir := filepath.Join(path, ".got")
	dirs := []string{"objects", "refs/heads", "refs/tags"}

	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(gotDir, d), 0755)
		if err != nil {
			return err
		}
	}

	head_path := filepath.Join(gotDir, "HEAD")

	err := os.WriteFile(head_path, []byte("ref: refs/heads/main\n"), 0644)
	if err != nil {
		return err
	}

	return nil
}
