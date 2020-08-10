package logic

import "testing"

func TestVerifyPassword(t *testing.T) {
	password := "1111111111111111111adfwerwersdfxzvgrtyryutuytuytfgdfhqwe11111111111111111111111111111"
	passwordErr := "222222"
	passwordDb, err := GeneratePassword(password)
	if err != nil {
		t.Log(err)
	}

	t.Log(len(passwordDb))
	if VerifyPassword(password, passwordDb) {
		t.Logf("%s success!", password)
	} else {
		t.Logf("%s failed", passwordErr)
	}

	if VerifyPassword(passwordErr, passwordDb) {
		t.Logf("%s success!", passwordErr)
	} else {
		t.Logf("%s failed", passwordErr)
	}
}
