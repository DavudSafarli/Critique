package main

import (
	"github.com/DavudSafarli/Critique/util/testing_utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testing_utils.PrepareTestDb()
	code := m.Run()
	os.Exit(code)
}
