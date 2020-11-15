package repository

import (
	"github.com/DavudSafarli/Critique/external/repository/abstract"
	"github.com/DavudSafarli/Critique/external/repository/postgres_repos"
	"strings"
)

// NewFeedbackRepository ..
func NewFeedbackRepository(driver string, connstr string) abstract.FeedbackRepository {
	switch strings.ToLower(driver) {
	case "pg":
		return postgres_repos.NewPGFeedbackRepository(connstr)
	}
	panic("NewFeedbackRepository not implemented for driver: " + driver)
}

func NewAttachmentRepository(driver string, connstr string) abstract.AttachmentRepository {
	switch strings.ToLower(driver) {
	case "pg":
		return postgres_repos.NewPGAttachmentRepository(connstr)
	}
	panic("NewFeedbackRepository not implemented for driver: " + driver)
}

// NewTagRepository creates new TagRepository
func NewTagRepository(driver string, connstr string) abstract.TagRepository {
	switch strings.ToLower(driver) {
	case "pg":
		return postgres_repos.NewPGTagRepository(connstr)
	}
	panic("NewTagRepository not implemented for driver: " + driver)
}
