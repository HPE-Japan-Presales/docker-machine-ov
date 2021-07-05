package driver

import (
	ov "github.com/HewlettPackard/oneview-golang/ov"
	log "github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/state"
)

type Oneview struct {
	Endpoint                  string `yaml:"endpoint"`
	ApiVersion                int    `yaml:"api-version"`
	Username                  string `yaml:"user"`
	Password                  string `yaml:"password"`
	Domain                    string `yaml:"domain,omitempty"`
	ServerProfileTemplateName string `yaml:"server-profile-template"`
	ServerProfileName         string `yaml:"server-profile"`
	ServerHardwareName        string `yaml:"server-hardware"`
}

// Precheck
func (o *Oneview) Validate() error {
	log.Debugf("OneView Structure: %#v", o)
	_, err := o.NewClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}

// Createt HPE OneView Client
func (o *Oneview) NewClient() (*ov.OVClient, error) {
	var ovClient *ov.OVClient

	//If pass invalid API ver, this will be panic.
	ovc := ovClient.NewOVClient(
		o.Username,
		o.Password,
		o.Domain,
		o.Endpoint,
		false, //ssl verificcation
		o.ApiVersion,
		"*")

	log.Debugf("Trying to connect HPE OneView endpoint at %v with API Ver %v", o.Endpoint, o.ApiVersion)
	_, err := ovc.GetAPIVersion()
	if err != nil {
		log.Error(Wrap(err))
		return nil, err
	}
	log.Debugf("Connected %v with API Ver %v", o.Endpoint, o.ApiVersion)

	return ovc, nil
}

func (o *Oneview) GetServerStatus() (string, error) {
	hardwareName := o.ServerHardwareName
	log.Infof("Get hardware status for %s", hardwareName)
	ovc, err := o.NewClient()
	if err != nil {
		log.Error(Wrap(err))
		return "", err
	}

	//Get Server hardware infomation
	hardware, err := ovc.GetServerHardwareByName(hardwareName)
	if err != nil {
		log.Error(Wrap(err))
		return "", err
	}
	return hardware.Status, nil
}

func (o *Oneview) PowerOn() error {
	ovc, err := o.NewClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	//Get Server hardware infomation
	hardwareName := o.ServerHardwareName
	hardware, err := ovc.GetServerHardwareByName(hardwareName)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	log.Infof("Power on %s", hardwareName)
	var powerTask *ov.PowerTask
	powerTask = powerTask.NewPowerTask(hardware)
	powerTask.GetCurrentPowerState()
	powerTask.PowerExecutor(1)
	return nil
}

func (o *Oneview) PowerOff() error {
	ovc, err := o.NewClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	//Get Server hardware infomation
	hardwareName := o.ServerHardwareName
	hardware, err := ovc.GetServerHardwareByName(hardwareName)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	log.Infof("Power off %s", hardwareName)
	var powerTask *ov.PowerTask
	powerTask = powerTask.NewPowerTask(hardware)
	powerTask.GetCurrentPowerState()
	powerTask.PowerExecutor(2)
	return nil
}

func (o *Oneview) GetPowerState() (state.State, error) {
	ovc, err := o.NewClient()
	if err != nil {
		log.Error(Wrap(err))
		return state.Error, err
	}

	//Get Server hardware infomation
	hardwareName := o.ServerHardwareName
	hardware, err := ovc.GetServerHardwareByName(hardwareName)
	if err != nil {
		log.Error(Wrap(err))
		return state.Error, err
	}

	log.Infof("Get server power state of %s", hardwareName)
	var powerTask *ov.PowerTask
	powerTask = powerTask.NewPowerTask(hardware)
	powerTask.GetCurrentPowerState()
	powerState := powerTask.State
	log.Debugf("Power state is %v of %s", powerState, hardwareName)

	switch powerState {
	case ov.P_ON:
		return state.Running, nil
	case ov.P_OFF:
		return state.Stopped, nil
	case ov.P_UKNOWN:
		return state.Error, nil
	default:
		return state.Error, nil
	}
	return state.Error, nil
}

func (o *Oneview) CreateServerProfile() error {
	ovc, err := o.NewClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	//Get Server Profile infomation
	serverProfileTemplateName := o.ServerProfileTemplateName
	serverProfileTemplate, err := ovc.GetProfileTemplateByName(serverProfileTemplateName)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	//Get Server hardware infomation
	hardwareName := o.ServerHardwareName
	hardware, err := ovc.GetServerHardwareByName(hardwareName)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	serverProfileName := o.ServerProfileName
	log.Infof("Create server profile %s from %s", serverProfileName, serverProfileTemplateName)

	if err := ovc.CreateProfileFromTemplate(serverProfileName, serverProfileTemplate, hardware); err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}

func (o *Oneview) DeleteServerProfile() error {
	ovc, err := o.NewClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	serverProfileName := o.ServerProfileName
	log.Infof("Delete server profile %s", serverProfileName)
	err = ovc.DeleteProfile(serverProfileName)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	//Wait delete completion
	return err
}
