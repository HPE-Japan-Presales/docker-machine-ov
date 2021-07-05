package driver

import (
	"testing"

	log "github.com/docker/machine/libmachine/log"
)

var (
	serverTestServer = &Server{
		//		Address:      "172.16.1.120",
		Address:      "172.16.14.10",
		RootPassword: "password",
		KsUrl:        "http://172.16.1.120/tak/ov-ks-label.iso",
		OsUrl:        "http://172.16.1.120/tak/CentOS-7-x86_64-Minimal-2003-ks.iso",
		SshPublicKey: "TEST_SSH_PUB_KEY",
	}
)

func TestServerValidate(t *testing.T) {
	log.SetDebug(true)
	server := serverTestServer
	if err := server.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestServerCopySshPubKey(t *testing.T) {
	log.SetDebug(true)
	server := serverTestServer
	if err := server.CopySshPubKey(); err != nil {
		t.Fatal(err)
	}
}

func TestServerWaitOsInstallation(t *testing.T) {
	log.SetDebug(true)
	server := serverTestServer
	//	server.Address = "172.16.14.10"
	if err := server.WaitOsInstallation(); err != nil {
		t.Fatal(err)
	}
}
