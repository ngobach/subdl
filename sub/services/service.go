package services

type Service interface {
	Search(kw string) []SubEntry
	Download(id string)
}
