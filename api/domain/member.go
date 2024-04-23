package domain

import "errors"

type MemberList []Member

func (memberList *MemberList) HasAdmin() bool {
	for _, member := range *memberList {
		if member.Role == Admin {
			return true
		}
	}

	return false
}

type Member struct {
	AccountId AccountId
	Role      Role
}

type Role int

const (
	Admin    = Role(iota)
	Customer = Role(iota)
)

var (
	ErrInvalidRole = errors.New("role_from_string: invalid role")
)

func RoleFromString(role string) (Role, error) {
	switch role {
	case "admin":
		return Admin, nil
	case "customer":
		return Customer, nil
	default:
		return Role(0), ErrInvalidRole
	}
}

func RoleToString(role Role) string {
	switch role {
	case Admin:
		return "admin"
	case Customer:
		return "customer"
	default:
		return ""
	}
}
