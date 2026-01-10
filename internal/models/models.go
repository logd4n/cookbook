package models

type Eat struct {
	ID           int      `json:"ID"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Ingredients  []string `json:"ingredients"`
	Instructions string   `json:"instructions"`
}

type RecipeShort struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
