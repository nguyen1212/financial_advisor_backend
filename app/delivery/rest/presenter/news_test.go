package presenter

import (
	"testing"

	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
)

func TestNew_FormNew(t *testing.T) {
	t.Parallel()
	assert.Equal(t, New{
		Title:     "title",
		Thumbnail: "thumbnail",
		Status:    "added",
		Content:   "finance",
	}, FormNew(dto.News{
		Title:     "title",
		Thumbnail: "thumbnail",
		Status:    "added",
		Content:   "finance",
	}))
}

func TestNew_FormNews(t *testing.T) {
	t.Parallel()
	newsDtos := []dto.News{
		{
			Title:     "title1",
			Thumbnail: "thumbnail1",
			Status:    "added",
			Content:   "finance1",
		},
		{
			Title:     "title2",
			Thumbnail: "thumbnail2",
			Status:    "synced",
			Content:   "finance2",
		},
	}
	expected := []New{
		{
			Title:     "title1",
			Thumbnail: "thumbnail1",
			Status:    "added",
			Content:   "finance1",
		},
		{
			Title:     "title2",
			Thumbnail: "thumbnail2",
			Status:    "synced",
			Content:   "finance2",
		},
	}

	assert.Equal(t, expected, FormNews(newsDtos))
}
