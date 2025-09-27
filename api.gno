package gnit

import (
	"gno.land/p/nt/avl"
)

func NewRepository(name string) *Repository {
	return &Repository{
		identity: Identity{
			Name: name,
		},
	}
}

func (r *Repository) Commit(message string, files map[string][]byte) string {
	if r.commits == nil {
		r.commits = avl.NewTree()
	}
	if r.objects == nil {
		r.objects = avl.NewTree()
	}

	treeHash := createTreeHash(files)

	tree := make(map[string]string) // path -> objectHash
	for path, content := range files {
		objectHash := createObjectHash(content)
		r.objects.Set(objectHash, content)
		tree[path] = objectHash
	}

	r.objects.Set(treeHash, tree)

	timestamp := int64(12345) // Simple fixed timestamp for now

	commit := &Commit{
		Tree:      treeHash,
		Parents:   []string{}, // TODO: implement parent tracking
		Author:    r.identity,
		Committer: r.identity,
		Message:   message,
		Timestamp: timestamp,
	}

	commitHash := createCommitHash(commit)
	commit.Hash = commitHash

	r.commits.Set(commitHash, commit)

	return commitHash
}

func (r *Repository) GetCommit(hash string) *Commit {
	if r.commits == nil {
		return nil
	}
	value, exists := r.commits.Get(hash)
	if !exists {
		return nil
	}
	return value.(*Commit)
}

func (r *Repository) GetFile(commitHash, path string) []byte {
	commit := r.GetCommit(commitHash)
	if commit == nil {
		return nil
	}

	treeValue, exists := r.objects.Get(commit.Tree)
	if !exists {
		return nil
	}

	tree := treeValue.(map[string]string)

	objectHash, exists := tree[path]
	if !exists {
		return nil
	}

	fileValue, exists := r.objects.Get(objectHash)
	if !exists {
		return nil
	}

	return fileValue.([]byte)
}
