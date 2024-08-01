package usecase

import "github.com/jackc/pgx/v5/pgtype"

type Waipoints struct {
	Id       string      `json:"id"`
	Sequence pgtype.Int4 `json:"sequence"`
}

type Results struct {
	Waipoints []Waipoints `json:"waypoints"`
}

type HereAPIRequest struct {
	Results []Results `json:"results"`
}
