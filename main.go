package tilde

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type currentUserLookupFunc func() (*user.User, error)
type otherUserLookupFunc func(string) (*user.User, error)

func expand(p string, currentUserLookup currentUserLookupFunc, otherUserLookup otherUserLookupFunc) (string, error) {
	l := strings.Split(p, string(os.PathSeparator))
	if len(l) < 1 {
		// it's empty so nothing we can do.  return p unchanged
		return p, nil
	}
	pfe := l[0]

	if l[0] == "" {
		// dunno what to do with it.  return p unchanged
		return p, nil
	} else if l[0] == "~" {
		// want homedir of current user
		if user, err := currentUserLookup(); err != nil {
			// but there was an error getting it.  return p unchanged
			return p, err
		} else {
			// all good.  set pfe for use at end.
			pfe = user.HomeDir
		}
	} else if l[0][0] != '~' {
		// it's not a homedir so nothing to do. return p unchanged
		return p, nil
	} else {
		// it's a homedir and that homedir is not the current one.
		// look it up and set pfe for use at the end
		if u, err := otherUserLookup(l[0][1:]); err != nil {
			// it failed though.  return p unchanged
			return p, err
		} else {
			pfe = u.HomeDir
		}
	}
	l[0] = pfe
	return filepath.Join(l...), nil
}

// Expand a tilde path with the user's homedir.
//
// Works for ~ and for ~user.
//
// In the event that a homedir can't be looked up, the original
// string is returned with an error.  The error may be discarded
// if the original string is what you want in the event of
// failure anyway.
//
// Does not support ~ at any other position other than the first char.
func Expand(p string) (string, error) {
	return expand(p, user.Current, user.Lookup)
}

// Try to expand.  Return original string on failure.
//
// This is similar to sh behavior.
func MayExpand(p string) string {
	exp, _ := expand(p, user.Current, user.Lookup)
	return exp
}

// Try to expand.  Panic on error.
//
// This is similar to zsh.
func MustExpand(p string) string {
	exp, err := expand(p, user.Current, user.Lookup)
	if err != nil {
		panic(err)
	}
	return exp
}
