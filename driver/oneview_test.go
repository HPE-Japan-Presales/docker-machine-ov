package driver

import (
	"fmt"
	"strconv"
	"testing"

	log "github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/state"
)

var (
	oneviewTestSimulator = map[string]string{
		"endpoint":                  "https://172.16.1.21",
		"apiVersion":                "1200",
		"username":                  "rancher",
		"password":                  "password",
		"serverProfileTemplateName": "Rancher-template",
		"serverProfileName":         "Rancher-test",
		"serverHardwareName":        "0000A66101, bay 5",
	}
)

func createTestOneview(oneviewType string) (*Oneview, error) {
	if oneviewType == "simulator" {
		apiVersion, err := strconv.Atoi(oneviewTestSimulator["apiVersion"])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse API version")
		}
		return &Oneview{
			Endpoint:                  oneviewTestSimulator["endpoint"],
			ApiVersion:                apiVersion,
			Username:                  oneviewTestSimulator["username"],
			Password:                  oneviewTestSimulator["password"],
			ServerProfileTemplateName: oneviewTestSimulator["serverProfileTemplateName"],
			ServerProfileName:         oneviewTestSimulator["serverProfileName"],
			ServerHardwareName:        oneviewTestSimulator["serverHardwareName"],
		}, nil
	}

	err := fmt.Errorf("%s is not supported oneview type", oneviewType)
	return nil, err
}

func TestOneviewPowerHandling(t *testing.T) {
	log.SetDebug(true)
	o, err := createTestOneview("simulator")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Power on %s", o.ServerHardwareName)
	if err := o.PowerOn(); err != nil {
		t.Fatal(err)
	}

	t.Logf("Get power state of %s", o.ServerHardwareName)
	powerState, err := o.GetPowerState()
	if err != nil {
		t.Fatal(err)
	}
	if powerState != state.Running {
		t.Fatal("Power state is not ON")
	}

	t.Logf("Power off %s", o.ServerHardwareName)
	if err := o.PowerOff(); err != nil {
		t.Fatal(err)
	}
	t.Logf("Get power state of %s", o.ServerHardwareName)
	powerState, err = o.GetPowerState()
	if err != nil {
		t.Fatal(err)
	}
	if powerState != state.Stopped {
		t.Fatal("Power state is not OFF")
	}
}

func TestOneviewServerProfileHandling(t *testing.T) {
	log.SetDebug(true)
	o, err := createTestOneview("simulator")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Create server profile %s from %s for %s", o.ServerProfileName, o.ServerProfileTemplateName, o.ServerHardwareName)
	if err := o.CreateServerProfile(); err != nil {
		t.Fatal(err)
	}

	t.Logf("Power on %s", o.ServerHardwareName)
	if err := o.PowerOn(); err != nil {
		t.Fatal(err)
	}

	t.Logf("Delete server profile %s", o.ServerProfileName)
	if err := o.DeleteServerProfile(); err != nil {
		t.Fatal(err)
	}

}
