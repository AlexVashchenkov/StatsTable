package models

type Context struct {
	Repository       string
	Revision         string
	ResolvedRevision string
	Extensions       []string
	Exclude          []string
	OrderBy          string

	RestrictTo   []string
	UseCommitter bool
	Languages    []string

	Format       string
	Files        []string
	BlameData    []BlameInfo
	Stats        []AuthorStats
	StatsToPrint []AuthorStatsForPrint
}

type BlameInfo struct {
	Author    string
	Committer string
	Lines     int
	File      string
}

type AuthorStats struct {
	Name    string
	Lines   int
	Commits map[string]struct{}
	Files   map[string]struct{}
}

type AuthorStatsForPrint struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}
