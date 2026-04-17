package heimdall

type resultset struct {
	Data       any   `json:"data,omitempty"`
	HasMore    bool  `json:"has_more,omitempty"`
	NextCursor int64 `json:"next_cursor,omitempty"`
}
