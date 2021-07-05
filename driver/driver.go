package driver

import (
	"fmt"
	"io/ioutil"
	"os"

	//	"strconv"
	"strings"

	"github.com/docker/machine/libmachine/drivers"
	log "github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

const (
	driverName     = "ov"
	defaultSshUser = "root"
	defaultSshPort = 22
)

type Driver struct {
	*drivers.BaseDriver
	*HpeConfig
}

func NewDriver(hostName, storePath string) *Driver {
	// DONT output anything message!! RPC server will be failed.
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
		HpeConfig: &HpeConfig{},
	}
}

// Create a host using the driver's config
func (d *Driver) Create() error {
	log.Info("Create server for HPE servers managed by HPE OneView")
	log.Debugf("BaseDriver: %#v", d.BaseDriver)
	log.Debugf("HpeConfig: %#v", d.HpeConfig)

	// Create server profile
	log.Info("Create server profile on HPE OneView")
	if err := d.HpeConfig.Oneview.CreateServerProfile(); err != nil {
		log.Error(Wrap(err))
		return err
	}

	// Create iLO client
	iloClient, err := d.HpeConfig.NewIloClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	// Insert virtual floppy for kickstart
	log.Info("Mount custom OS image on HPE iLO virtual DVD")
	if err := iloClient.InsertVirtualMedia(d.HpeConfig.Server.OsUrl, "dvd"); err != nil {
		log.Error(Wrap(err))
		return err
	}
	defer iloClient.EjectVirtualMedia("dvd")

	// Insert virtual floppy for kickstart
	log.Info("Mount kickstart image on HPE iLO virtual Floppy")
	if err := iloClient.InsertVirtualMedia(d.HpeConfig.Server.KsUrl, "floppy"); err != nil {
		log.Error(Wrap(err))
		return err
	}
	defer iloClient.EjectVirtualMedia("floppy")

	// Power on to install OS
	log.Info("Power on server")
	if err := d.HpeConfig.Oneview.PowerOn(); err != nil {
		log.Error(Wrap(err))
		return err
	}

	// Wait OS install
	log.Info("Start OS installation")
	d.HpeConfig.Server.WaitOsInstallation()

	//Prepare ssh key pair
	log.Info("Create ssh keys")
	if err := d.genSshKeyPairs(); err != nil {
		log.Error(Wrap(err))
		return err
	}

	// Copy public key
	log.Info("Copy ssh keys")
	if err := d.HpeConfig.Server.CopySshPubKey(); err != nil {
		log.Error(Wrap(err))
		return err
	}

	log.Info("Server setup has been done!")
	return nil
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return driverName
}

// GetCreateFlags returns the mcnflag.Flag slice representing the flags
// that can be set, their descriptions and defaults.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return mcnFlags
}

// GetIP returns an IP or hostname that this host is available at
// e.g. 1.2.3.4 or docker-host-d60b70a14d3a.cloudapp.net
func (d *Driver) GetIP() (string, error) {
	if d.BaseDriver.IPAddress == "" {
		err := fmt.Errorf("Could not get IP Address: %#v", d.BaseDriver)
		log.Error(Wrap(err))
		return "", err
	}
	return d.BaseDriver.IPAddress, nil
}

// GetMachineName returns the name of the machine
func (d *Driver) GetMachineName() string {
	return d.BaseDriver.MachineName
}

// GetSSHHostname returns hostname for use with ssh
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) genSshKeyPairs() error {
	// Create ssh key pair for new server
	sshPrivateKeyPath := d.GetSSHKeyPath()
	sshPublicKeyPath := fmt.Sprintf("%s.pub", d.GetSSHKeyPath())
	log.Infof("Ssh private key path: %s", sshPrivateKeyPath)
	log.Infof("Ssh public key path: %s", sshPublicKeyPath)

	if err := ssh.GenerateSSHKey(sshPrivateKeyPath); err != nil {
		log.Error(Wrap(err))
		return err
	}
	sshPrivateKey, err := ioutil.ReadFile(sshPrivateKeyPath)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	sshPublicKey, err := ioutil.ReadFile(sshPublicKeyPath)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	log.Infof("Ssh private key: %s", sshPrivateKey)
	log.Infof("Ssh public key: %s", sshPublicKey)
	d.HpeConfig.Server.SshPrivateKey = string(sshPrivateKey)
	d.HpeConfig.Server.SshPublicKey = strings.TrimSuffix(string(sshPublicKey), "\n")

	return nil
}

// GetSSHPort returns port for use with ssh
func (d *Driver) GetSSHPort() (int, error) {
	if d.BaseDriver.SSHPort == 0 {
		err := fmt.Errorf("Could not get SSH port: %#v", d.BaseDriver)
		log.Error(Wrap(err))
		return 0, err
	}
	return d.BaseDriver.SSHPort, nil
}

// GetSSHUsername returns username for use with ssh
func (d *Driver) GetSSHUsername() string {
	return d.BaseDriver.SSHUser
}

// GetURL returns a Docker compatible host URL for connecting to this host
// e.g. tcp://1.2.3.4:2376
func (d *Driver) GetURL() (string, error) {
	ip := d.HpeConfig.Server.Address
	if ip == "" {
		err := fmt.Errorf("Could not get Docker IP address: %#v", d.HpeConfig)
		log.Error(Wrap(err))
		return "", err
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

// GetState returns the state that the host is in (running, stopped, etc)
func (d *Driver) GetState() (state.State, error) {
	powerState, err := d.HpeConfig.Oneview.GetPowerState()
	if err != nil {
		log.Error(Wrap(err))
		return state.Error, err
	}
	return powerState, nil
}

// PreCreateCheck allows for pre-create operations to make sure a driver is ready for creation
func (d *Driver) PreCreateCheck() error {
	log.Info("Check HPE OneView configurations")
	err := d.HpeConfig.Oneview.Validate()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	log.Info("Check new server configurations")
	err = d.HpeConfig.Server.Validate()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	return nil
}

// Remove a host
func (d *Driver) Remove() error {
	log.Info("Start to remove server profile from HPE OneView")
	if err := d.HpeConfig.Oneview.DeleteServerProfile(); err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}

// SetConfigFromFlags configures the driver with the object that was returned
// by RegisterCreateFlags
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	// Debug flag
	if os.Getenv(strings.ToUpper(driverName)+"_DEBUG") != "" {
		log.Info("Turn on debug mode from env value")
		log.SetDebug(true)
	}
	if flags.Bool(driverName + "-debug") {
		log.Info("Turn on debug mode from command option")
		log.SetDebug(true)
	}

	// Read parameters
	if flags.String(driverName+"-yaml") != "" {
		// From yaml file
		log.Info("Configuration is read from yaml")
		d.HpeConfig.Yaml = &Yaml{
			Path: flags.String(driverName + "-yaml"),
		}
		conf, err := d.HpeConfig.Yaml.Read()
		if err != nil {
			log.Error(Wrap(err))
			return err
		}
		d.HpeConfig = conf
	} else {
		// From command option
		log.Info("Configuration is read from command options")
		d.HpeConfig = &HpeConfig{
			Oneview: &Oneview{
				Endpoint:                  flags.String(driverName + "-oneview-endpoint"),
				ApiVersion:                flags.Int(driverName + "-oneview-api-version"),
				Username:                  flags.String(driverName + "-oneview-user"),
				Password:                  flags.String(driverName + "-oneview-password"),
				Domain:                    flags.String(driverName + "-oneview-domain"),
				ServerProfileTemplateName: flags.String(driverName + "-oneview-server-profile-template"),
				ServerHardwareName:        flags.String(driverName + "-oneview-server-hardware"),
			},
			Server: &Server{
				Address:      flags.String(driverName + "-server-address"),
				RootPassword: flags.String(driverName + "-server-root-password"),
				KsBaseUrl:    flags.String(driverName + "-server-kickstart-base-url"),
				KsUrl:        fmt.Sprintf("%s/%s.iso", flags.String(driverName+"-server-kickstart-base-url"), flags.String(driverName+"-server-address")),
				OsUrl:        flags.String(driverName + "-server-os-url"),
			},
		}
	}

	d.BaseDriver.IPAddress = d.HpeConfig.Server.Address
	d.BaseDriver.SSHUser = defaultSshUser
	d.BaseDriver.SSHPort = defaultSshPort
	d.HpeConfig.Server.Hostname = d.GetMachineName()
	d.HpeConfig.Oneview.ServerProfileName = fmt.Sprintf("%s-docker-machine-%s", driverName, d.GetMachineName())
	d.HpeConfig.Server.KsUrl = fmt.Sprintf("%s/%s.iso", d.HpeConfig.Server.KsBaseUrl, d.HpeConfig.Server.Address)

	log.Debugf("BaseDriver: %#v", d.BaseDriver)
	log.Debugf("HpeConfig: %#v", d.HpeConfig)

	return nil
}

// Start a host
func (d *Driver) Start() error {
	if err := d.HpeConfig.Oneview.PowerOn(); err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}

// Stop a host gracefully
func (d *Driver) Stop() error {
	if err := d.HpeConfig.Oneview.PowerOff(); err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}

// Restart a host. This may just call Stop(); Start() if the provider does not
// have any special restart behaviour.
func (d *Driver) Restart() error {
	err := d.Stop()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	err = d.Start()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}

// Kill stops a host forcefully
func (d *Driver) Kill() error {
	err := d.Stop()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}
