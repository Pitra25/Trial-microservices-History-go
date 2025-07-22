package types

type Recording struct {
	ID          int16  `json:"id"`
	Calculation string `json:"calculation"`
	CreatedAt   string `json:"createdAt"`
}

type BodyStructure struct {
	Calculation string `json:"Calculation"`
	CreatedAt   string `json:"CreatedAt"`
}
