package integrationtest

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"testing"

	"pyg.com/api/app"
	"pyg.com/api/app/data"
	apperrors "pyg.com/api/app/errors"
	"pyg.com/api/domain"
	"pyg.com/api/infra"
)

func arrangeDatabase() (*sql.DB, data.SessionToken) {
	db := infra.ConnectToDatabase()
	cleanDb(db)
	
	return db, createAccountAndLogin(db)
}

func createAccountAndLogin(db *sql.DB) data.SessionToken {
	return createAccountWithTagAndLogin("", db)
}

func createAccountWithTagAndLogin(tag string, db *sql.DB) data.SessionToken {
	email := fmt.Sprintf("test%s@gmail.com", tag)
	username := fmt.Sprintf("username_%s", tag)

	err := app.Register(email, username, "FirstName", "LastName", "Pass", db)
	if err != nil { log.Fatalln(err) }

	sessionToken, err := app.Login(username, "Pass", db)
	if err != nil { log.Fatalln(err) }

	return sessionToken
}

func TestCreateDeveloper(t *testing.T) {
	db, sessionToken := arrangeDatabase()
	defer db.Close()

	_, err := app.CreateDeveloper(sessionToken, "DeveloperName", db)
	if err != nil { t.Fatal(err) }
}

func TestUpdateDeveloper(t *testing.T) {
	db, sessionToken := arrangeDatabase()
	defer db.Close()

	developerId, err := app.CreateDeveloper(sessionToken, "DeveloperName", db)
	if err != nil { t.Fatal(err) }

	newMember := []domain.Member {
		{
			AccountId: domain.AccountId(0),
			Role: domain.Admin,
		},
	}
	newName := "NewDeveloperName"
	err = app.UpdateDeveloper(sessionToken, developerId, &newName, newMember, db)
	if err != nil { t.Fatal("failed to update developer") }

	updatedDeveloper, err := app.GetDeveloper(developerId, db)
	if err != nil { t.Fatal(err) }

	if updatedDeveloper.Name != domain.DeveloperName(newName) { t.Fatalf("unexpected name: %s", updatedDeveloper.Name) }
	if !slices.Equal(updatedDeveloper.Members, newMember) {
		t.Fatalf("unexpected member list") 
	}
}

func TestUpdateDeveloperFailsWithoutPermission(t *testing.T) {
	db := infra.ConnectToDatabase()
	cleanDb(db)
	defer db.Close()

	sessionA := createAccountWithTagAndLogin("A", db)
	sessionB := createAccountWithTagAndLogin("B", db)

	developerId, err := app.CreateDeveloper(sessionA, "developer", db)
	if err != nil { t.Fatal(err) }

	newDeveloper := "new developer"
	err = app.UpdateDeveloper(sessionB, developerId, &newDeveloper, nil, db)
	if !errors.Is(err, apperrors.ErrUpdateUnauthorized) {
		t.Fatal("update should have failed")
	}
}

func TestGetDevelopers(t *testing.T) {
	db, sessionToken := arrangeDatabase()
	defer db.Close()

	for i := 0; i < 4; i++ {
		_, err := app.CreateDeveloper(sessionToken, fmt.Sprintf("%d", i), db)
		if err != nil { log.Fatalln(err) }
	}

	pageToken := data.PageToken(0)

	developers, pageToken, err := app.GetDevelopers(pageToken, 2, db)
	if err != nil { t.Fatal(err) }
	if !containsAllListed(developers, 0, 1) { t.Fatalf("only developers 0 and 1 were expected")}

	developers, pageToken, err = app.GetDevelopers(pageToken, 2, db)
	if err != nil { t.Fatal(err) }
	if !containsAllListed(developers, 2, 3) { t.Fatalf("only developers 2 and 3 were expected")}
	
	developers, _, err = app.GetDevelopers(pageToken, 2, db)
	if err != nil { t.Fatal(err) }
	if !containsAllListed(developers) { t.Fatalf("no developers were expected")}
}

func containsAllListed(developers []domain.Developer, ids ...int) bool {
	if len(developers) != len(ids) {
		return false
	}

	for _, developer := range developers {
		id, err := strconv.Atoi(string(developer.Name))
		if err != nil { panic("unexpected conversion failure") }

		if !slices.Contains(ids, id) {
			return false
		}
		
	}

	return true
}
