package git

type RepositoryManager struct {
	repos map[string]Repository
}

// Repository is a repository cloned and managed for use across reconciliations.
// TODO: decide what information we need to hold in a "Repository".
type Repository struct {
	Path string
	URL  string
}
