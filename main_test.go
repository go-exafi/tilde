package tilde

import (
	"os/user"
	"testing"
)

func currentUserLookupMock() (*user.User, error) {
	return otherUserLookupMock("testuser")
}

func otherUserLookupMock(u string) (*user.User, error) {
	if u == "other" {
		return &user.User{
			Uid:      "124",
			Gid:      "457",
			Username: "otheruser",
			Name:     "Other User",
			HomeDir:  "/home/other",
		}, nil
	}
	if u == "other user" {
		return &user.User{
			Uid:      "125",
			Gid:      "458",
			Username: "other user",
			Name:     "Other User2",
			HomeDir:  "/home/other user",
		}, nil
	}
	if u == "testuser" {
		return &user.User{
			Uid:      "123",
			Gid:      "456",
			Username: "testuser",
			Name:     "Test User",
			HomeDir:  "/home/test",
		}, nil
	}
	return nil, user.UnknownUserError("could not find user")
}

type TestCase struct {
	tildePath    string
	expectedPath string
	err          bool
}

func TestExpand_current(t *testing.T) {
	tests := []TestCase{
		{"~", "/home/test", false},
		{"~other", "/home/other", false},
		{"~other/.ssh/known_hosts", "/home/other/.ssh/known_hosts", false},
		{"~other user/.ssh/known_hosts", "/home/other user/.ssh/known_hosts", false},
		{"~other test user/.ssh/known_hosts", "", true},
		{"~george", "", true},
	}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		t.Run(test.tildePath, func(t *testing.T) {
			home, err := expand(test.tildePath, currentUserLookupMock, otherUserLookupMock)
			if err != nil && !test.err {
				t.Errorf("Error looking up ~: %v", err)
			} else if _, ok := err.(user.UnknownUserError); test.err && !ok {
				t.Errorf("Expected error not found: %v", test.err)
			} else if !test.err && home != test.expectedPath {
				t.Errorf("Returned path %#v was not expected path %#v", home, test.expectedPath)
			}
		})
	}
}
