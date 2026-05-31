package endpoint

import "strings"

type Path string

func NewPath(path string) (Path, error) {
	trimmedPath := strings.TrimSpace(path)
	if len(trimmedPath) == 0 {
		return "", ErrEmptyPath
	}
	if strings.ContainsRune(trimmedPath, ' ') {
		return "", ErrPathHasSpaces
	}
	return Path(trimmedPath), nil
}
