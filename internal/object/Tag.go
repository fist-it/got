package object

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

type Tag struct {
	Object  [32]byte
	ObjType string
	Name    string
	Tagger  Author
	Message string
}

func (t *Tag) Type() string {
	return "tag"
}

// tag <size>\0
// object <hex>\ntype <type>\ntag <name>\ntagger <name> <email> <timestamp> <tz>\n\n<message>
func (t *Tag) Serialize() []byte {
	var result []byte

	var body []byte
	body = append(body, fmt.Sprintf("object %x\n", t.Object)...)
	body = append(body, fmt.Sprintf("type %s\n", t.ObjType)...)
	body = append(body, fmt.Sprintf("tag %s\n", t.Name)...)
	body = append(body, fmt.Sprintf("tagger %s\n", t.Tagger.String())...)

	body = append(body, '\n')
	body = append(body, fmt.Sprintf("%s", t.Message)...)

	result = append(result, fmt.Sprintf("tag %d\x00", len(body))...)
	result = append(result, body...)

	return result
}

func Deserialize(data []byte) (*Tag, error) {
	before, after, _ := bytes.Cut(data, []byte{0})
	header, _, _ := bytes.Cut(before, []byte{' '})
	body := after

	if string(header) != "tag" {
		return nil, errors.New("Invalid tag header")
	}

	lines := bytes.Split(body, []byte("\n"))
	// field: object
	i := 0
	if !bytes.HasPrefix(lines[i], []byte("object")) {
		return nil, errors.New("Invalid object header")
	}
	object_hash := lines[i][7:]
	decoded_object, err := hex.DecodeString(string(object_hash))
	if err != nil {
		return nil, errors.New("Error parsing object hash")
	}
	object := [32]byte(decoded_object)
	i++

	// field: type
	if !bytes.HasPrefix(lines[i], []byte("type")) {
		return nil, errors.New("Invalid type header")
	}
	obj_type := string(lines[i][5:])
	i++

	// field: tag
	if !bytes.HasPrefix(lines[i], []byte("tag")) {
		return nil, errors.New("Invalid tag name header")
	}
	tag_name := string(lines[i][4:])
	i++

	// field: tagger
	if !bytes.HasPrefix(lines[i], []byte("tagger")) {
		return nil, errors.New("Invalid tagger header")
	}
	author_to_parse := lines[i][7:]
	var author *Author
	author, err = ParseAuthor(string(author_to_parse))
	if err != nil {
		return nil, errors.New("Error parsing author string")
	}
	i++

	// empty line before the tag message
	i++
	message := string(bytes.Join(lines[i:], []byte("\n")))

	return &Tag{
		Object:  object,
		ObjType: obj_type,
		Name:    tag_name,
		Tagger:  *author,
		Message: message,
	}, nil
}
