package driver

import (
	"testing"

	log "github.com/docker/machine/libmachine/log"
)

var (
	iloTestOneview = &Oneview{
		Endpoint:                  "https://192.168.2.6",
		ApiVersion:                1200,
		Username:                  "rancher",
		Password:                  "password",
		ServerProfileTemplateName: "Rancher-template",
		ServerProfileName:         "Rancher-test",
		ServerHardwareName:        "SGH652SV73, bay 5",
	}
	iloTestServer = &Server{
		KsUrl: "http://172.16.1.120/tak/ov-ks-label.iso",
		OsUrl: "http://172.16.1.120/tak/CentOS-7-x86_64-Minimal-2003-ks.iso",
	}
)

var hpeConfig *HpeConfig

func createTestIloClient() (*IloClient, error) {
	hpeConfig = &HpeConfig{
		Oneview: iloTestOneview,
		Server:  iloTestServer,
	}
	client, err := hpeConfig.NewIloClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func TestIloVirtualMediaActions(t *testing.T) {
	log.SetDebug(true)
	c, err := createTestIloClient()
	if err != nil {
		t.Fatal(err)
	}

	dvdImageUrl := hpeConfig.Server.OsUrl
	floppyImageUrl := hpeConfig.Server.KsUrl
	log.Infof("Image URL: %s", dvdImageUrl)

	err = c.InsertVirtualMedia(dvdImageUrl, "dvd")
	if err != nil {
		t.Fatal(err)
	}
	err = c.InsertVirtualMedia(floppyImageUrl, "floppy")
	if err != nil {
		t.Fatal(err)
	}
	err = c.EjectVirtualMedia("dvd")
	if err != nil {
		t.Fatal(err)
	}
	err = c.EjectVirtualMedia("floppy")
	if err != nil {
		t.Fatal(err)
	}
}
