package sln

type Solution struct {
	Projects []Project
}

type Project struct {
	ID          string
	Name        string
	ProjectFile string
	TypeGUID    string
}
