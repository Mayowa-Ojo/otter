package internal

import "testing"

const (
	succeed = "\u2713"
	failed  = "\u2717"
	actk    = "007235f4-d5a0-4d75-b354-b20a65e6b87a"
	rftk    = "e773c02c-15d8-42c5-8b96-ar3a27364700"
)

func TestPersistAuthorization(t *testing.T) {
	t.Log("Should save auth tokens to config file in home directory")
	{
		if err := PersistAuthorization(actk, rftk); err != nil {
			t.Fatalf("\t%s\tShould create directory if not exist: %v", failed, err)
		}

		t.Logf("\t%s\tShould create directory if not exist: ", succeed)
	}
}

func TestVerifyAuthToken(t *testing.T) {
	t.Log("Should return false for an invalid token")
	{
		token := "007235f4-d5a0-4d75-s354-b20a65e6b87v"

		isValid := VerifyAuthToken(token)

		if !isValid {
			t.Logf("\t%s\ttoken is invalid", succeed)
		} else {
			t.Fatalf("\t%s\ttoken is valid", failed)
		}
	}
}
