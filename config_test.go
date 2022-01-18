package bx24sync

import (
	"testing"
)

func Test_Socket_String_Success(t *testing.T) {

	sock := Socket{Host: "127.0.0.1", Port: 9000}

	if "127.0.0.1:9000" != sock.String() {
		t.Errorf("Expected %s. Result %s")
	}

}

func Test_Endpoint_Url_Success(t *testing.T) {
	t.Fail()
}

func Test_BasicAuth_GetPair_Success(t *testing.T) {
	t.Fail()
}

func Test_NewRegistrarConfigFromEnv_Success(t *testing.T) {
	t.Fail()
}

func Test_NewDataBaseConnectionFromEnv_Success(t *testing.T) {
	t.Fail()
}

func Test_DataBaseConnection_MakeConnURL_Success(t *testing.T) {
	t.Fail()
}

func Test_NewConsumerConfigFromEnv_Success(t *testing.T) {
	t.Fail()
}
