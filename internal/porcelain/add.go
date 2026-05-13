package porcelain

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fist-it/got/internal/index"
	"github.com/fist-it/got/internal/object"
	"github.com/fist-it/got/internal/repo"
	"github.com/fist-it/got/internal/store"
)

func modeFor(info os.FileInfo) uint32 {
	m := info.Mode()
	switch {
	case m&os.ModeSymlink != 0:
		return 0120000
	case m.IsDir():
		return 0040000 // no dir should reach add
	case m&0111 != 0:
		return 0100755 // executable
	default:
		return 0100644 // regular
	}
}

func Add(path string) error {
	repoRoot, err := repo.Find()
	if err != nil {
		return err
	}

	got_dir := filepath.Join(repoRoot, ".got")
	str := &store.Store{Root: filepath.Join(got_dir, "objects")}
	idx, err := index.Load(got_dir)
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	err = filepath.WalkDir(absPath, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if d.Name() == ".got" {
				return fs.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}

		info, err := os.Lstat(p)
		if err != nil {
			return err
		}

		st, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return errors.New("error reading Stat_t of file")
		}
		relPath, err := filepath.Rel(repoRoot, p)
		if err != nil {
			return err
		}
		if strings.HasPrefix(relPath, "..") {
			return errors.New("Error constructing relative path: leading up a directory")
		}
		relPath = filepath.ToSlash(relPath)

		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		blob := &object.Blob{Content: data}
		hash, err := str.Write(blob)
		if err != nil {
			return err
		}

		entry := &index.Entry{
			CTime: time.Unix(st.Ctim.Sec, st.Ctim.Nsec),
			MTime: info.ModTime(),
			Dev:   uint32(st.Dev),
			Ino:   uint32(st.Ino),
			Mode:  modeFor(info),
			Uid:   st.Uid,
			Gid:   st.Gid,
			Size:  uint32(info.Size()),
			Hash:  hash,
			Flags: 0,
			Path:  relPath,
		}
		idx.Add(entry)

		return nil
	})
	if err != nil {
		return err
	}

	err = idx.Save(got_dir)
	if err != nil {
		return err
	}

	return nil
}
