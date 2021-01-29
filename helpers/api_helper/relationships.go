package api_helper

type RelationshipsObject struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type RelationshipsResponse struct {
	Data RelationshipsObject `json:"data"`
}

func (r *RelationshipsObject) JSONAPIFormat() map[string]RelationshipsResponse {
	relationships := make(map[string]RelationshipsResponse)
	relationships[r.Type] = RelationshipsResponse{
		Data: *r,
	}

	return relationships
}
