package gnit

import "strconv"

func createTreeHash(files map[string][]byte) string {
	var content string

	var keys []string
	for k := range files {
		keys = append(keys, k)
	}

	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	for _, filename := range keys {
		content += filename + string(files[filename])
	}

	return simpleHash(content)
}

func createObjectHash(content []byte) string {
	return simpleHash(string(content))
}

func createCommitHash(commit *Commit) string {
	content := commit.Tree + commit.Message + commit.Author.Name + commit.Author.Email + strconv.FormatInt(commit.Timestamp, 10)

	for _, parent := range commit.Parents {
		content += parent
	}

	return simpleHash(content)
}

func simpleHash(input string) string {
	var hash int64 = 5381

	for i := 0; i < len(input); i++ {
		hash = ((hash << 5) + hash) + int64(input[i])
	}

	if hash < 0 {
		hash = -hash
	}

	return strconv.FormatInt(hash, 16)
}
