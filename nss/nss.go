package main

import (
	. "github.com/protosam/go-libnss"
	. "github.com/protosam/go-libnss/structs"
)

// Placeholder main() stub is neccessary for compile.
func main() {}

func init() {
	// We set our implementation to "TestImpl", so that go-libnss will use the methods we create
	SetImpl(TestImpl{})
}

// We're creating a struct that implements LIBNSS stub methods.
type TestImpl struct{ LIBNSS }

////////////////////////////////////////////////////////////////
// Passwd Methods
////////////////////////////////////////////////////////////////

// PasswdAll() will populate all entries for libnss
func (self TestImpl) PasswdAll() (Status, []Passwd) {
	if len(dbtest_passwd) == 0 {
		return StatusUnavail, []Passwd{}
	}
	return StatusSuccess, dbtest_passwd
}

// PasswdByName() returns a single entry by name.
func (self TestImpl) PasswdByName(name string) (Status, Passwd) {
	for _, entry := range dbtest_passwd {
		if entry.Username == name {
			return StatusSuccess, entry
		}
	}
	return StatusNotfound, Passwd{}
}

// PasswdByUid() returns a single entry by uid.
func (self TestImpl) PasswdByUid(uid uint) (Status, Passwd) {
	for _, entry := range dbtest_passwd {
		if entry.UID == uid {
			return StatusSuccess, entry
		}
	}
	return StatusNotfound, Passwd{}
}

////////////////////////////////////////////////////////////////
// Group Methods
////////////////////////////////////////////////////////////////
// endgrent
func (self TestImpl) GroupAll() (Status, []Group) {
	if len(dbtest_group) == 0 {
		return StatusUnavail, []Group{}
	}
	return StatusSuccess, dbtest_group
}

// getgrent
func (self TestImpl) GroupByName(name string) (Status, Group) {
	for _, entry := range dbtest_group {
		if entry.Groupname == name {
			return StatusSuccess, entry
		}
	}
	return StatusNotfound, Group{}
}

// getgrnam
func (self TestImpl) GroupByGid(gid uint) (Status, Group) {
	for _, entry := range dbtest_group {
		if entry.GID == gid {
			return StatusSuccess, entry
		}
	}
	return StatusNotfound, Group{}
}

////////////////////////////////////////////////////////////////
// Shadow Methods
////////////////////////////////////////////////////////////////
// endspent
func (self TestImpl) ShadowAll() (Status, []Shadow) {
	if len(dbtest_shadow) == 0 {
		return StatusUnavail, []Shadow{}
	}
	return StatusSuccess, dbtest_shadow
}

// getspent
func (self TestImpl) ShadowByName(name string) (Status, Shadow) {
	for _, entry := range dbtest_shadow {
		if entry.Username == name {
			return StatusSuccess, entry
		}
	}
	return StatusNotfound, Shadow{}
}
