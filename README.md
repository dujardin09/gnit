## Overview
Composable Git protocol implementation as a Gno package (`gno.land/p/demo/gnit`) that enables any realm to become a Git repository.

## Core Principles
* **Protocol-only**: Git functionality without built-in permissions
* **Composable**: Import into any realm to add Git capabilities
* **Safe Object**: Repository objects can be exposed directly with internal safety

## API Specification
### Core Types
```go
package gnit

type Repository struct { /* opaque */ }

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

type TreeEntry struct {
    Mode FileMode
    Name string
    Hash string
}
```
### Repository API
```go
// Constructor
func NewRepository(name string) *Repository

// Write operations
func (r *Repository) Commit(message string, files map[string][]byte) string
func (r *Repository) CreateBranch(name, from string) error
func (r *Repository) DeleteBranch(name string) error
func (r *Repository) Merge(target, source string) (string, error)
func (r *Repository) Tag(name, commit string) error

// Read operations
func (r *Repository) GetCommit(hash string) *Commit
func (r *Repository) GetHistory(ref string, limit int) []Commit
func (r *Repository) GetFile(commit, path string) []byte
func (r *Repository) GetTree(commit, path string) []TreeEntry
func (r *Repository) GetBranches() []string
func (r *Repository) GetTags() map[string]string
func (r *Repository) GetDiff(from, to string) Diff

// Paginated operations for large data transfers
func (r *Repository) GetHistoryPaged(ref string, page, pageSize int) ([]Commit, bool)
func (r *Repository) GetTreePaged(commit, path string, page, pageSize int) ([]TreeEntry, bool)
func (r *Repository) GetFilePaged(commit, path string, offset, chunkSize int) ([]byte, bool)
func (r *Repository) GetPackfile(from, to string, page int) (PackChunk, bool)
func (r *Repository) GetPackInfo(from, to string) PackInfo  // Pagination metadata

// XXX: Permissions
func (r *Repository) SetPermissionChecker(func(op string, caller std.Address) bool)
```
### Pagination Types
```go
type PackChunk struct {
    Page     int
    Total    int
    Objects  []PackObject
    HasMore  bool
}

type PackObject struct {
    Hash    string
    Type    ObjectType
    Size    int
    Data    []byte
}
```
## Usage Patterns
### Direct Exposure
```go
package myproject

import "gno.land/p/gnit"

var Repository *gnit.Repository

func init() {
    Repository = gnit.NewRepository("myproject")
}
// Users call: myproject.Repository.Commit("fix", files)
```
### Custom Access Control
```go
package teamproject

import (
    "gno.land/p/gnit"
    "gno.land/p/demo/avl"
)

var (
    repo    *gnit.Repository
    members *avl.Tree
)

func init() {
    repo = gnit.NewRepository("teamproject")
    repo.SetPermissionChecker(checkPermission)
}

func checkPermission(op string, caller std.Address) bool {
    role, _ := members.Get(caller.String())
    return role == "developer" || role == "admin"
}

func Commit(message string, files map[string][]byte) string {
    return repo.Commit(message, files)
}
```
### Multi-Repository
```go
package hub

import "gno.land/p/gnit"

var repos *avl.Tree  // name -> *gnit.Repository

func CreateRepo(name string) *gnit.Repository {
    repo := gnit.NewRepository(name)
    repos.Set(name, repo)
    return repo
}

func GetRepo(name string) *gnit.Repository {
    r, _ := repos.Get(name)
    return r.(*gnit.Repository)
}
```
## CLI Integration
# Configure remote
gnit remote add origin gno.land/r/user/project

# Basic operations
```bash
gnit push                   # → Repository.Commit()
gnit pull                   # → Repository.GetHistoryPaged() (multiple calls)
gnit clone <realm>          # → Repository.GetPackfile() (auto-resumes using cache)
gnit branch <name>          # → Repository.CreateBranch()
gnit merge <branch>         # → Repository.Merge()
```
### Clone/Pull Protocol
The CLI uses `.gnit/objects/` to cache fetched objects and automatically continues from where it left off:

```go
// Repository API additions for pagination info
func (r *Repository) GetPackInfo(from, to string) PackInfo

type PackInfo struct {
    TotalObjects int
    TotalPages   int
    PageSize     int
    CommitHash   string
}
```
### Read Constraints
Blockchain read operations have process limits, requiring:

* **Chunked transfers** for large repositories
* **Pagination** for history and file listings
* **Progressive loading** during clone/pull
* **Configurable chunk sizes** based on network limits

## Security Model
1. **Default**: Only owning realm can write
2. **Custom**: Via SetPermissionChecker()
3. **Wrapped**: Realm implements custom methods

## Extensions
Future considerations:

* Large files (IPFS)
* Signed commits
* Pull requests (separate protocol)
* Submodules 🤮
* Hooks 🤮

## MVP Specification
### Minimal Git Objects
* **Commit**: Hash, message, tree reference, parent, timestamp
* **Tree**: Directory listing with file/subdirectory entries
* **Blob**: File content storage
* **Single branch**: `main` only, no branching operations

### Repository Operations
```go
// Minimal API surface
func NewRepository(name string) *Repository
func (r *Repository) Commit(message string, files map[string][]byte) string
func (r *Repository) GetHistory(limit int) []Commit
func (r *Repository) GetFile(commit, path string) []byte
func (r *Repository) GetPackInfo() PackInfo
func (r *Repository) GetPackfile(page int) PackChunk
```
### CLI Commands
* **`gnit clone <realm>`**: Fetch repository with automatic resume
* **`gnit log`**: Display commit history
* **`gnit show <commit>`**: Show commit details and files

### Local Storage
* **`.gnit/objects/`**: Content-addressed object cache
* **`.gnit/config`**: Remote realm configuration
* **Automatic resume**: Re-running clone continues from cached objects

### Limitations
* Linear history only (fast-forward commits)
* Single branch support
* No merge commits (no divergent history)
* Read-only operations initially
