package repository

import (
	"testing"
	pb_account "github.com/onezerobinary/db-box/proto/account"
	"fmt"
	"github.com/goinggo/tracelog"
)

func TestConnectToDB(t *testing.T){

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	db, err := ConnectToDB()

	if err != nil  {
		t.Error("Wrong DB or error in connection: ", db)
	}

	tracelog.Completed("testDB", "TestConnectToDB")
}

func TestAddDoc(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	//fakeToken := pb_account.Token{"fff"}
	fakeStatus := pb_account.Status{pb_account.Status_DISABLED}

	username := "Pino"
	password := "Mancuso"

	faketoken := GenerateToken(username, password)

	fakeAccount := pb_account.Account{
		faketoken,
		username,
		password,
		nil,
		&fakeStatus,
		"Account",
		"2018-01-11",
		"2028-01-10",
	}

	token, err :=  AddDoc(fakeAccount)

	if err != nil {
		t.Error("non possible to store the document into the DB")
	}

	fmt.Println("ok, token generated: ", token)
}

func TestIsPresent(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	uuid := "2e9dbe901ff0172692e2e38e39f8daf5ae2f7c71"

	p := IsPresent(uuid)

	if !p {
		t.Error("Element not present: ", uuid)
	}
}

func TestGetAccountByCredentials(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	fakeCredentials := pb_account.Credentials{}

	fakeCredentials.Username = "johnny"
	fakeCredentials.Password = "doe"

	token := pb_account.Token{}

	to := GenerateToken(fakeCredentials.Username, fakeCredentials.Password)

	token.Token = to
	fakeCredentials.Token = &token

	account, err := GetAccountByCredentials(fakeCredentials)

	if err != nil {
		t.Error("Error in getting the informations", err)
	}

	if account == nil {
		t.Errorf("Error, no account found")
	}
}

func TestGetAccountByToken(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	fakeCredentials := pb_account.Credentials{}

	fakeCredentials.Username = "johnny"
	fakeCredentials.Password = "doe"

	token := pb_account.Token{}

	to := GenerateToken(fakeCredentials.Username, fakeCredentials.Password)

	token.Token = to

	account, err := GetAccountByToken(token)

	if err != nil {
		t.Error("Error in getting the informations", err)
	}

	if account == nil {
		t.Errorf("Error, no account found")
	}
}

func TestRemoveDoc(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	fakeCredentials := pb_account.Credentials{}

	fakeCredentials.Username = "john"
	fakeCredentials.Password = "doe"

	token := pb_account.Token{}

	to := GenerateToken(fakeCredentials.Username, fakeCredentials.Password)

	token.Token = to

	err := RemoveDoc(token)

	if err != nil {
		t.Errorf("Error to delete the account")
	}
}

func TestUpdateDoc(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	// Create and add an account to the DB
	fakeStatus := pb_account.Status{pb_account.Status_DISABLED}

	username := "Mary"
	password := "Rossi"

	faketoken := GenerateToken(username, password)

	fakeAccount := pb_account.Account{
		faketoken,
		username,
		password,
		nil,
		&fakeStatus,
		"Account",
		"2018-01-11",
		"2028-01-10",
	}

	token, err :=  AddDoc(fakeAccount)

	tk := pb_account.Token{Token:token}

	account, err := GetAccountByToken(tk)

	account.Token.Token = token
	//change the something to the account
	account.Status = &pb_account.Status{pb_account.Status_SUSPENDED}

	err = UpdateDoc(*account)

	if err != nil {
		t.Errorf("Error to update the account")
	}
}

func TestCheckEmail(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	email := pb_account.Email{ "Mary"}

	token, err := CheckEmail(email)

	if err != nil {
		t.Errorf("Error to get the email back")
	}

	if token != nil {
		response := "Token: " + token.Token
		tracelog.Trace("", "", response)
	}
}

func TestSetAccountStatus(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	token := pb_account.Token{Token:"49926336a43fc66a762e64fa5b1bf1a86fa20848"}
	account, err := GetAccountByToken(token)

	if err != nil {
		tracelog.Errorf(err,"database_test", "TestSetAccountStatus", "Account not found")
	}

	if account != nil {

		account.Status = &pb_account.Status{pb_account.Status_REVOKED}
		accountUpdateStatus := pb_account.UpdateStatus{&token, account.Status}

		err = SetAccountStatus(accountUpdateStatus)

		if err != nil {
			tracelog.Errorf(err,"database_test", "TestSetAccountStatus", "Account status not updated")
			t.Errorf("Status not updated")
		}
	}
}

func TestGetAccountStatus(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	token := pb_account.Token{Token:"49926336a43fc66a762e64fa5b1bf1a86fa20848"}
	account, err := GetAccountByToken(token)

	if err != nil {
		tracelog.Errorf(err,"database_test", "TestGetAccountStatus", "Account not found")
	}

	if account != nil {
		accountStatus, err := GetAccountStatus(token)

		if err != nil {
			tracelog.Errorf(err,"database_test", "TestGetAccountStatus", "Account Status not found")
		}

		if accountStatus != nil {
			text := "Account Status retrieved: " + accountStatus.Status.String()
			tracelog.Trace("database_test", "TestGetAccountStatus" , text)
		} else {
			t.Errorf("Account status not retrieved")
		}
	}
}


func TestGetAccountsByStatus(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	accountStatus := pb_account.Status{pb_account.Status_ENABLED}

	acconts, err := GetAccountsByStatus(accountStatus)

	if err != nil {
		tracelog.Errorf(err,"database_test", "TestGetAccountsByStatus", "Accounts not found")
	}

	if acconts != nil {
		text := "Acounts retrieved with status: " + accountStatus.Status.String()
		tracelog.Trace("database_test", "TestGetAccountsByStatus" , text)
	} else {
		t.Errorf("Accounts not retrieved")
	}

}

func TestGetAccounts(t *testing.T) {

	tracelog.Start(tracelog.LevelTrace)
	defer tracelog.Stop()

	acconts, err := GetAccounts()

	if err != nil {
		tracelog.Errorf(err,"database_test", "TestGetAccounts", "Accounts not found")
	}

	if acconts != nil {
		text := "Accounts retrieved"
		tracelog.Trace("database_test", "TestGetAccounts" , text)
	} else {
		t.Errorf("Accounts not retrieved")
	}
}