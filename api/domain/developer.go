package domain

import "strings"

type DeveloperId int64

type Developer struct {
	Id        DeveloperId
	Name      string
	CreatedBy AccountId
	Members   MemberList
}

func MakeDeveloper(name string, createdBy AccountId) *Developer {
	return &Developer{
		Name:      name,
		CreatedBy: createdBy,
		Members: MemberList{{
			AccountId: createdBy,
			Role:      Admin,
		}},
	}
}

func (developer *Developer) IsValid() bool {
	return developer.Members.HasAdmin() && strings.TrimSpace(developer.Name) != ""
}

func (developer *Developer) IsAdmin(accountId AccountId) bool {
	for _, member := range developer.Members {
		if member.AccountId == accountId {
			return true
		}
	}
	return false
}
