package gnit

import "testing"

func TestNewRepository(t *testing.T) {
	r := NewRepository("foo")

	if r.identity.Name != "foo" {
		t.Errorf("expected %s, got %s", "foo", r.identity.Name)
	}
}

func TestCommit(t *testing.T) {
	r := NewRepository("test-repo")

	files := map[string][]byte{
		"main.go":   []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"),
		"README.md": []byte("# Test Project\n\nThis is a test project."),
	}

	commitHash := r.Commit("Initial commit", files)

	// Verify commit hash is generated
	if commitHash == "" {
		t.Error("expected non-empty commit hash")
	}

	// Verify commit is stored
	commit := r.GetCommit(commitHash)
	if commit == nil {
		t.Error("expected commit to be stored")
	}

	if commit.Message != "Initial commit" {
		t.Errorf("expected message 'Initial commit', got %s", commit.Message)
	}

	if commit.Hash != commitHash {
		t.Errorf("expected hash %s, got %s", commitHash, commit.Hash)
	}
}

func TestGetFile(t *testing.T) {
	r := NewRepository("test-repo")

	files := map[string][]byte{
		"main.go":   []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"),
		"README.md": []byte("# Test Project\n\nThis is a test project."),
	}

	commitHash := r.Commit("Initial commit", files)

	// Test retrieving existing files
	mainContent := r.GetFile(commitHash, "main.go")
	if string(mainContent) != string(files["main.go"]) {
		t.Errorf("expected %s, got %s", string(files["main.go"]), string(mainContent))
	}

	readmeContent := r.GetFile(commitHash, "README.md")
	if string(readmeContent) != string(files["README.md"]) {
		t.Errorf("expected %s, got %s", string(files["README.md"]), string(readmeContent))
	}

	// Test retrieving non-existing file
	nonExistentContent := r.GetFile(commitHash, "nonexistent.txt")
	if nonExistentContent != nil {
		t.Error("expected nil for non-existent file")
	}
}

func TestGetCommitNonExistent(t *testing.T) {
	r := NewRepository("test-repo")

	commit := r.GetCommit("nonexistent-hash")
	if commit != nil {
		t.Error("expected nil for non-existent commit")
	}
}

func TestMultipleCommits(t *testing.T) {
	r := NewRepository("test-repo")

	// First commit
	files1 := map[string][]byte{
		"file1.txt": []byte("First version"),
	}
	hash1 := r.Commit("First commit", files1)

	// Second commit
	files2 := map[string][]byte{
		"file1.txt": []byte("Second version"),
		"file2.txt": []byte("New file"),
	}
	hash2 := r.Commit("Second commit", files2)

	// Verify both commits exist and are different
	if hash1 == hash2 {
		t.Error("expected different hashes for different commits")
	}

	// Verify first commit content
	content1 := r.GetFile(hash1, "file1.txt")
	if string(content1) != "First version" {
		t.Errorf("expected 'First version', got %s", string(content1))
	}

	// Verify second commit content
	content2 := r.GetFile(hash2, "file1.txt")
	if string(content2) != "Second version" {
		t.Errorf("expected 'Second version', got %s", string(content2))
	}

	newFile := r.GetFile(hash2, "file2.txt")
	if string(newFile) != "New file" {
		t.Errorf("expected 'New file', got %s", string(newFile))
	}
}
