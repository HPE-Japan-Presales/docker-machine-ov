package driver

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/docker/machine/libmachine/drivers"
)

const (
	driverTestDefaultHostname = "bay5"
	//	driverTestDefaultHostname                         = "bay4"
	driverTestDefaultStorePath                        = "/tmp"
	driverTestDefaultOneviewEndpoint                  = "https://192.168.2.6"
	driverTestDefaultOneviewApiVersion                = 1200
	driverTestDefaultOneviewUsername                  = "rancher"
	driverTestDefaultOneviewPassword                  = "password"
	driverTestDefaultOneviewServerProfileTemplateName = "Rancher-template"
	driverTestDefaultOneviewServerHardwareName        = "SGH652SV73, bay 5"
	//	driverTestDefaultOneviewServerHardwareName = "SGH652SV73, bay 4"
	driverTestDefaultServerAddress          = "172.16.14.10"
	driverTestDefaultServerRootPassword     = "password"
	driverTestDefaultServerKickstartBaseUrl = "http://172.16.1.120/tak"
	driverTestDefaultServerOsUrl            = "http://172.16.1.120/tak/CentOS-7-x86_64-Minimal-2003-ks.iso"
)

func createTestDriver() (*Driver, error) {
	driver := NewDriver(driverTestDefaultHostname, driverTestDefaultStorePath)
	flags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			fmt.Sprintf(driverName + "-oneview-endpoint"):                driverTestDefaultOneviewEndpoint,
			fmt.Sprintf(driverName + "-oneview-api-version"):             driverTestDefaultOneviewApiVersion,
			fmt.Sprintf(driverName + "-oneview-user"):                    driverTestDefaultOneviewUsername,
			fmt.Sprintf(driverName + "-oneview-password"):                driverTestDefaultOneviewPassword,
			fmt.Sprintf(driverName + "-oneview-server-profile-template"): driverTestDefaultOneviewServerProfileTemplateName,
			fmt.Sprintf(driverName + "-oneview-server-hardware"):         driverTestDefaultOneviewServerHardwareName,
			fmt.Sprintf(driverName + "-server-address"):                  driverTestDefaultServerAddress,
			fmt.Sprintf(driverName + "-server-root-password"):            driverTestDefaultServerRootPassword,
			fmt.Sprintf(driverName + "-server-kickstart-base-url"):       driverTestDefaultServerKickstartBaseUrl,
			fmt.Sprintf(driverName + "-server-os-url"):                   driverTestDefaultServerOsUrl,
			fmt.Sprintf(driverName + "-debug"):                           true,
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	if err := driver.SetConfigFromFlags(flags); err != nil {
		return nil, err
	}
	return driver, nil
}

func TestDriverCreate(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Create(); err != nil {
		t.Fatal(err)
	}
}

func TestDriverGetSSHKeyPath(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	sshKeyPath := d.GetSSHKeyPath()
	if sshKeyPath == "" {
		t.Fatal("could not get ssh key path")
	}
	t.Log(sshKeyPath)
}

func TestDriverDriverName(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	driverName := d.DriverName()
	if driverName == "" {
		t.Fatalf("Failed to get Driver Name: %s", driverName)
	}
	t.Log(driverName)
}

func TestDriverGetCreateFlags(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	flags := d.GetCreateFlags()
	if len(flags) == 0 {
		t.Fatalf("Flags is not set: %s", flags)
	}
	t.Log(flags)
}

func TestDriverGetIP(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	ipAddress, err := d.GetIP()
	if err != nil {
		t.Fatal(err)
	}
	if ipAddress == "" {
		t.Fatalf("Not handling invalid IP Address: %s", ipAddress)
	}
	t.Logf(ipAddress)
}

func TestDriverGetMachineName(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	name := d.GetMachineName()
	if name == "" {
		t.Fatalf("Not handling invalid Machine Name: %s", name)
	}
	t.Logf("Machine Name: %s", name)
}

func TestDriverGetSSHHostname(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	name, err := d.GetSSHHostname()
	if err != nil {
		t.Fatal(err)
	}
	if name == "" {
		t.Fatalf("Not handling invalid SSH Host name: %s", name)
	}
	t.Logf("SSH Host name: %s", name)
}

func TestDriverGetSSHPort(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	port, err := d.GetSSHPort()
	if err != nil {
		t.Fatal(err)
	}
	if port == 0 {
		t.Fatalf("Not handling invalid SSH port: %v\n", port)
	}
	t.Logf("SSH Port: %v\n", port)
}

func TestDriverGetSSHUsername(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	name := d.GetSSHUsername()
	if name == "" {
		t.Fatalf("Not handling invalid SSH User Name: %s\n", name)
	}
	t.Logf("SSH User Name: %s\n", name)
}

func TestDriverGetURL(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	url, err := d.GetURL()
	if err != nil {
		t.Fatal(err)
	}
	if url == "" {
		t.Fatalf("Not handling invalid URL: %s\n", url)
	}
	t.Logf("Docker URL: %s\n", url)
}

func TestDriverGetState(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	state, err := d.GetState()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Docker Status: %v\n", state)
}

func TestDriverKill(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Kill(); err != nil {
		t.Fatal(err)
	}
	state, err := d.GetState()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Docker Status: %v\n", state)
}

func TestDriverRemove(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Remove(); err != nil {
		t.Fatal(err)
	}
}

func TestDriverRestart(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Restart(); err != nil {
		t.Fatal(err)
	}
}

func TestDriverStart(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Start(); err != nil {
		t.Fatal(err)
	}
}

func TestDriverStop(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestDriverGenSshKeyPairs(t *testing.T) {
	d, err := createTestDriver()
	if err != nil {
		t.Fatal(err)
	}
	path := strings.Replace(d.GetSSHKeyPath(), "/id_rsa", "", -1)
	t.Logf("PATH: %s", path)

	if err := os.MkdirAll(path, 0777); err != nil {
		t.Fatal(err)
	}
	if err := d.genSshKeyPairs(); err != nil {
		t.Fatal(err)
	}
	d.HpeConfig.Server.Address = "172.16.1.120"
	if err := d.HpeConfig.Server.CopySshPubKey(); err != nil {
		t.Fatal(err)
	}
}
