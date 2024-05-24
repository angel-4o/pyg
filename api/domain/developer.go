package domain

// import "strings"

type DeveloperId int64
type DeveloperName string


type Developer struct {
	Id        DeveloperId
	Name      DeveloperName
	CreatedBy AccountId
	Members   MemberList
}

func MakeDeveloper(name string, createdBy AccountId) *Developer {
	return &Developer{
		Name:      DeveloperName(name),
		CreatedBy: createdBy,
		Members: MemberList{{
			AccountId: createdBy,
			Role:      Admin,
		}},
	}
}

func (developer *Developer) IsValid() bool {
	return developer.Members.HasAdmin() && developer.Name != ""
}

func (developer *Developer) IsAdmin(accountId AccountId) bool {
	for _, member := range developer.Members {
		if member.AccountId == accountId {
			return true
		}
	}
	return false
}
