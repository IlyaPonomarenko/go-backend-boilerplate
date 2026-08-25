package model

// Post is the resource returned by the API. JSON tags define the public field
// names used in request and response bodies.
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// PostInput contains fields clients may send when creating or replacing a
// post. The server owns the ID, so clients cannot choose it in the body.
type PostInput struct {
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
