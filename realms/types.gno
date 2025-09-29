package gnit

import (
	"std"

	"gno.land/p/nt/avl"
)

type Repository struct {
	identity Identity

	head    string    // current branch name
	refs    *avl.Tree // branch (string) -> hash
	commits *avl.Tree // hash -> []byte
	objects *avl.Tree // hash -> Object
}

type Commit struct {
	Hash      string
	Tree      string
	Parents   []string
	Author    Identity
	Committer Identity
	Message   string
	Timestamp int64
}

type Identity struct {
	Name    string
	Email   string
	Address std.Address
}

// type TreeEntry struct {
// 	Mode FileMode
// 	Name string
// 	Hash string
// }
