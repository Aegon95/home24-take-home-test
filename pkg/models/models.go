package models

type WebStats struct {
	HTMLVersion string
	Title       string

	H1Count int
	H2Count int
	H3Count int
	H4Count int
	H5Count int
	H6Count int

	InternalLinksCount int
	ExternalLinksCount int

	InaccessibleLinksCount int32

	HasLogin bool
}

