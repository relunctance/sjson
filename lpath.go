package sjson

import "strings"

type lpath struct {
	left  string
	right string
}

func newLpath(path string) *lpath {
	l := &lpath{}
	pos := strings.Index(path, "#")
	if pos == -1 {
		return l
	}
	l.left = strings.Trim(path[:pos], ".")
	l.right = strings.Trim(path[pos+1:], ".")
	return l
}
func (l *lpath) IsFinish() bool {
	return strings.Index(l.right, "#") == -1
}

func (l *lpath) String() string {
	s := ""
	s += "left: " + l.left
	s += " right: " + l.right
	s += "\n"
	return s
}
