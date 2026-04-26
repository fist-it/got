package object

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

type Commit struct {
	Tree      [32]byte   // Tree hash
	Parents   [][32]byte // Parent commits
	Committer Author
	Author    Author
	Message   string
}

func (c *Commit) Type() string {
	return "commit"
}

func (c *Commit) Serialize() []byte {
	var result []byte

	var body []byte
	body = append(body, fmt.Sprintf("tree %x\n", c.Tree)...)
	for _, parent := range c.Parents {
		body = append(body, fmt.Sprintf("parent %x\n", parent)...)
	}

	body = append(body, fmt.Sprintf("author %s\n", c.Author.String())...)
	body = append(body, fmt.Sprintf("committer %s\n", c.Committer.String())...)
	body = append(body, '\n')
	body = append(body, []byte(c.Message)...)

	result = fmt.Appendf(nil, "commit %d\x00", len(body))
	result = append(result, body...)

	return result
}

func DeserializeCommit(data []byte) (*Commit, error) {

	before, after, _ := bytes.Cut(data, []byte{0})
	header := before
	body := after

	before_space, _, _ := bytes.Cut(header, []byte{' '})
	commit_header := before_space
	if string(commit_header) != "commit" {
		return nil, errors.New("Invalid commit header")
	}

	lines := bytes.Split(body, []byte("\n"))

	i := 0
	tree_header := string(lines[i])[0:5]
	if string(tree_header) != "tree " {
		return nil, errors.New("Invalid tree header")
	}
	tree, err := hex.DecodeString(string(lines[i])[5:])

	if err != nil {
		return nil, errors.New("Error parsing the tree hash")
	}
	i++

	var parents [][32]byte
	for bytes.HasPrefix(lines[i], []byte("parent")) {
		parent_hash := lines[i][7:]
		decoded_parent, err := hex.DecodeString(string(parent_hash))
		if err != nil {
			return nil, errors.New("Error parsing a parent hash")
		}
		parents = append(parents, [32]byte(decoded_parent))
		i++
	}

	author_header := string(lines[i])
	space_idx := bytes.IndexByte([]byte(author_header), ' ')
	if string(author_header[0:space_idx]) != "author" {
		return nil, errors.New("Error parsing the author header: ")
	}

	author, err := ParseAuthor(string(lines[i]))
	if err != nil {
		return nil, errors.New("Error parsing the author: ")
	}
	i++

	committer, err := ParseAuthor(string(lines[i]))
	if err != nil {
		return nil, errors.New("Error parsing the committer: ")
	}
	i++

	// empty line
	i++
	message := string(bytes.Join(lines[i:], []byte("\n")))

	return &Commit{
		Tree:      [32]byte(tree),
		Parents:   parents,
		Committer: *committer,
		Author:    *author,
		Message:   message,
	}, nil
}
