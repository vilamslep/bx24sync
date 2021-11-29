package registrator_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/vi-la-muerto/bx24-service/rest/registrator"
	"github.com/vi-la-muerto/bx24-service/scheme"
)

func TestServiceClientsSuccess(t *testing.T) {
	
	srv := registrator.NewServer(15000)
	go srv.Run()
	
	defer srv.Close()

	client := scheme.Client{
		Name: "TEST_NAME",
		Code: "TEST_CODE",
		UID:  "TEST_UID",
	}

	content := strings.NewReader(client.String())

	res, err := http.Post(fmt.Sprintf("http://localhost:%d/client", srv.Port), "application/json", content)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatal(err)
	}

	
	defer res.Body.Close()

}
