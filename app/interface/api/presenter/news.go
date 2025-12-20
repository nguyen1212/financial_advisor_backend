// Package presenter handle logic to transform data returned to client side.
package presenter

type New struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Status    string `json:"status"`
	Content   string `json:"content,omitempty"`
}

func FormNew() New {
	return New{
		Title:     "Sample News Title",
		Thumbnail: "https://example.com/thumbnail.jpg",
		Status:    "published",
	}
}

func FormNews() []New {
	return []New{
		FormNew(),
		{
			Title:  "Another News Title",
			Status: "draft",
		},
	}
}
