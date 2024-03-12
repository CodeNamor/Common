package path

import (
	"path"
	"unicode/utf8"
)

const (
	slash = '/'
)

// Resolve converts a sequence of path segments into an resolved path.
// If any of the segments is absolute then the result will be a
// cleaned absolute path, otherwise it will be a cleaned join of the
// path segments.
// Modeled after Node.js path.resolve it processes from right to left
// looking for any absolute paths, then it joins that and any
// subsequent paths. Makes it easy to find path from a starting dir.
// Using path rather than filepath since we are dealing with unix
// paths not windows (and forward slashes work too).
func Resolve(paths ...string) string {
	segments := make([]string, 0, len(paths))
	for _, segment := range paths {
		if len(segment) == 0 { // empty string, ignore
			continue
		}
		c, _ := utf8.DecodeRuneInString(segment) // get first rune of string
		if c == slash {                          // segment starts with a slash, absolute path
			// since we have an absolute path, clear out any previous segments
			// this will be the staring point now
			segments = segments[:0] // clear the slice, retain capacity
		}
		segments = append(segments, segment)
	}

	// join the remaining segments and clean the result
	return path.Join(segments...)
}
