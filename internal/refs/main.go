package refs

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

type Refs struct {
	GotDir string
}

func (r *Refs) ResolveHEAD() ([32]byte, error) {
	head_path := filepath.Join(r.GotDir, "HEAD")

	head_data, err := os.ReadFile(head_path)
	if err != nil {
		return [32]byte{}, err
	}

	content := strings.TrimSpace(string(head_data))

	if strings.HasPrefix(content, "ref: ") {
		refPath := content[5:]
		return r.ReadRef(refPath)
	} else {
		hash, err := hex.DecodeString(content)
		if err != nil {
			return [32]byte{}, err
		}
		return [32]byte(hash), nil
	}
}

func (r *Refs) ReadRef(ref_path string) ([32]byte, error) {
	full_path := filepath.Join(r.GotDir, ref_path)
	data, err := os.ReadFile(full_path)
	if err != nil {
		return [32]byte{}, err
	}
	hash, err := hex.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		return [32]byte{}, err
	}
	return [32]byte(hash), nil
}

func (r *Refs) UpdateRef(ref_path string, hash [32]byte) error {
	fullPath := filepath.Join(r.GotDir, ref_path)
	content := hex.EncodeToString(hash[:]) + "\n"
	return os.WriteFile(fullPath, []byte(content), 0644)
}

func (r *Refs) CreateRef(name string, currentCommitHash [32]byte) error {
	ref_path := "refs/heads/" + name
	err := os.MkdirAll(filepath.Join(r.GotDir, "refs/heads/"), 0755)
	if err != nil {
		return err
	}
	return r.UpdateRef(ref_path, currentCommitHash)
}

func (r *Refs) ListRefs() (map[string][32]byte, error) {
	result := make(map[string][32]byte)
	refsDir := filepath.Join(r.GotDir, "refs")

	err := filepath.WalkDir(refsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(r.GotDir, path)
		if err != nil {
			return err
		}

		hash, err := r.ReadRef(relPath)
		if err != nil {
			return err
		}

		result[filepath.ToSlash(relPath)] = hash
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
