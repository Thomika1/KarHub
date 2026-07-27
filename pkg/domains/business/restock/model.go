package restock

type PriorityItem struct {
	PartID         string `json:"partId"`
	Name           string `json:"name"`
	CurrentStock   *int   `json:"currentStock,omitempty"`
	ProjectedStock *int   `json:"projectedStock,omitempty"`
	MinimumStock   *int   `json:"minimumStock,omitempty"`
	UrgencyScore   *int   `json:"urgencyScore,omitempty"`
}

type PrioritiesResponse struct {
	Priorities []PriorityItem `json:"priorities"`
}
