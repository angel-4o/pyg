package auth

import (
	"database/sql"
	"slices"

	"pyg.com/api/app/data"
	"pyg.com/api/domain"
	"pyg.com/api/persistence"
)

func HasAdminRole(sessionToken data.SessionToken, members domain.MemberList, db *sql.DB) bool {
	sessionDataSource := persistence.MakeSessionDataSource(db)
	session, err := sessionDataSource.FindSessionByToken(sessionToken)
	if err != nil { return false }

	memberIndex := slices.IndexFunc(members, func(member domain.Member) bool {
		return member.AccountId == session.AccountId && member.Role == domain.Admin
	})

	return memberIndex != -1
}