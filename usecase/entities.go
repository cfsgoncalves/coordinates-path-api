package usecase

import "github.com/jackc/pgx/v5/pgtype"

type Waipoints struct {
	Id       pgtype.Int4
	Sequence int
}

type Results struct {
	Waipoints []Waipoints
}

type HereAPIRequest struct {
	Results []Results
}
