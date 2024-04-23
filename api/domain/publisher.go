package domain

import "strings"

type PublisherId int64

type Publisher struct {
	Id        PublisherId
	CreatedBy AccountId
	Name      string
	Members   MemberList
}

func MakePublisher(name string, createdBy AccountId) *Publisher {
	return &Publisher{
		Name:      name,
		CreatedBy: createdBy,
		Members: MemberList{{
			AccountId: createdBy,
			Role:      Admin,
		}},
	}
}

func (publisher *Publisher) IsValid() bool {
	return publisher.Members.HasAdmin() && strings.TrimSpace(publisher.Name) != ""
}

func (publisher *Publisher) IsAdmin(accountId AccountId) bool {
	for _, member := range publisher.Members {
		if member.AccountId == accountId {
			return true
		}
	}
	return false
}
