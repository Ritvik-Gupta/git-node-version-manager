package utils

type Repository struct {
	Name, Url string
}

type RepoWithPackagesFound struct {
	Repository
	Packages []Package
}
