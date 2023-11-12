package repository

const (
	AscDirection  = "ASC"
	DescDirection = "DESC"
)

type SortDto struct {
	By        string
	Direction string
}
