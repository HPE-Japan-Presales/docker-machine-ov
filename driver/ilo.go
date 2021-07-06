package driver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/HewlettPackard/oneview-golang/rest"
	log "github.com/docker/machine/libmachine/log"
	"github.com/stmcginnis/gofish"
)

type IloClient struct {
	Address          string
	RemoteConsoleUrl string `json:"remoteConsoleUrl"`
	Token            string
	Model            string
	//	VirtualMedia     *IloVirtualMedia
	VirtualDevices *VirtualDevices
}

type VirtualDevices struct {
	Dvd    IloVirtualMedia
	Floppy IloVirtualMedia
}

type IloVirtualMedias struct {
	Members      []IloVirtualMediaMember
	MembersCount int `json:"Members@odata.count"`
}

type IloVirtualMedia struct {
	Id         string                 `json:"@odata.id"`
	MediaTypes []string               `json:"MediaTypes"`
	Actions    IloVirtualMediaActions `json:"Actions"`
}

type IloVirtualMediaActions struct {
	Insert IloVirtualMediaActionTarget
	Eject  IloVirtualMediaActionTarget
}

type Ilo5VirtualMedia struct {
	Id         string                  `json:"@odata.id"`
	MediaTypes []string                `json:"MediaTypes"`
	Actions    Ilo5VirtualMediaActions `json:"Actions"`
}

type Ilo5VirtualMediaActions struct {
	Insert IloVirtualMediaActionTarget `json:"#VirtualMedia.InsertMedia"`
	Eject  IloVirtualMediaActionTarget `json:"#VirtualMedia.EjectMedia"`
}

type Ilo4VirtualMedia struct {
	Id         string              `json:"@odata.id"`
	MediaTypes []string            `json:"MediaTypes"`
	Oem        Ilo4VirtualMediaOem `json:"Oem"`
}

type Ilo4VirtualMediaOem struct {
	Hp Ilo4VirtualMediaOemHp `json:"Hp"`
}

type Ilo4VirtualMediaOemHp struct {
	Actions Ilo4VirtualMediaActions `json:"Actions"`
}

type Ilo4VirtualMediaActions struct {
	Insert IloVirtualMediaActionTarget `json:"#HpiLOVirtualMedia.InsertVirtualMedia"`
	Eject  IloVirtualMediaActionTarget `json:"#HpiLOVirtualMedia.EjectVirtualMedia"`
}

type IloVirtualMediaActionTarget struct {
	Target string `json:"target"`
}

type IloVirtualMediaMember struct {
	Id string `json:"@odata.id"`
}

type IloInsertVirtualMediaReqBody struct {
	Image string `json:"Image"`
}

func (s *HpeConfig) NewIloClient() (*IloClient, error) {
	log.Info("Create new HPE iLO client")
	log.Debugf("HpeConfig: %#v", s)

	ovc, err := s.Oneview.NewClient()
	log.Debugf("OneviewClient: %#v", ovc)
	if err != nil {
		log.Error(Wrap(err))
		return nil, err
	}

	//Get Server hardware infomation
	hardwareName := s.Oneview.ServerHardwareName
	hardware, err := ovc.GetServerHardwareByName(hardwareName)
	log.Debugf("Hardware: %#v", hardware)
	if err != nil {
		log.Error(Wrap(err))
		return nil, err
	}

	// Retrieve HPE iLO session token from HPE OneView
	uri := fmt.Sprintf("%v/remoteConsoleUrl", hardware.URI)
	ovc.RefreshLogin()
	remoteConsoleResp, err := ovc.RestAPICall(rest.GET, uri, nil)
	log.Debugf("remoteConsoleResp: %#v", remoteConsoleResp)
	if err != nil {
		log.Error(Wrap(err))
		return nil, err
	}
	var iloClient IloClient
	if err := json.Unmarshal(remoteConsoleResp, &iloClient); err != nil {
		log.Error(Wrap(err))
		return nil, err
	}
	// Set Values
	iloClient.Address = hardware.MpHostInfo.MpIPAddresses[1].Address
	iloClient.Token = strings.Replace(iloClient.RemoteConsoleUrl, fmt.Sprintf("hplocons://addr=%v&sessionkey=", iloClient.Address), "", -1)
	iloClient.Model = hardware.MpModel
	log.Debugf("iloClient: %#v", iloClient)

	if iloClient.Address == "" {
		err := fmt.Errorf("Could not retrieve HPE iLO address")
		log.Error(Wrap(err))
		return nil, err
	}
	if iloClient.Token == "" {
		err := fmt.Errorf("Could not retrieve HPE iLO token")
		log.Error(Wrap(err))
		return nil, err
	}

	return &iloClient, nil
}

func (ilo *IloClient) GetVirtualMedia() error {
	log.Info("Get iLO Virtual Media infomation")

	// Create RedFish client
	c, err := ilo.createRedfishClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	defer c.Logout()

	// Retrieve managers on HPE iLO
	service := c.Service
	manager, err := service.Managers()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	// Retrieve virtual media members on HPE iLO
	// Gofish seems not support virtual media...
	resVirtualMedia, err := c.Get(fmt.Sprintf("%vVirtualMedia/", manager[0].Entity.ODataID))
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	defer resVirtualMedia.Body.Close()

	var virtualMedias IloVirtualMedias
	body, err := ioutil.ReadAll(resVirtualMedia.Body)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	if err := json.Unmarshal(body, &virtualMedias); err != nil {
		log.Error(Wrap(err))
		return err
	}

	virtualDevices := &VirtualDevices{}
	for _, virtualMediaMember := range virtualMedias.Members {
		res, err := c.Get(virtualMediaMember.Id)
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(Wrap(err))
			return err
		}

		iloModel := ilo.Model
		switch iloModel {
		case "iLO4":
			var virtualMedia Ilo4VirtualMedia
			if err := json.Unmarshal(body, &virtualMedia); err != nil {
				log.Error(Wrap(err))
				return err
			}

			for _, mediaType := range virtualMedia.MediaTypes {
				if mediaType == "DVD" {
					log.Info("iLO4 DVD virtual media detected")
					log.Debugf("Virtual DVD: %#v", virtualMedia)
					virtualDevices.Dvd = IloVirtualMedia{
						Id:         virtualMedia.Id,
						MediaTypes: virtualMedia.MediaTypes,
						Actions: IloVirtualMediaActions{
							Insert: virtualMedia.Oem.Hp.Actions.Insert,
							Eject:  virtualMedia.Oem.Hp.Actions.Eject,
						},
					}
				}
				if mediaType == "Floppy" {
					log.Info("iLO4 Floppy virtual media detected")
					log.Debugf("Virtual Floppy: %#v", virtualMedia)
					virtualDevices.Floppy = IloVirtualMedia{
						Id:         virtualMedia.Id,
						MediaTypes: virtualMedia.MediaTypes,
						Actions: IloVirtualMediaActions{
							Insert: virtualMedia.Oem.Hp.Actions.Insert,
							Eject:  virtualMedia.Oem.Hp.Actions.Eject,
						},
					}
				}
			}
			ilo.VirtualDevices = virtualDevices

		case "iLO5":
			var virtualMedia Ilo5VirtualMedia
			if err := json.Unmarshal(body, &virtualMedia); err != nil {
				log.Error(Wrap(err))
				return err
			}
			log.Debugf("%v", string(body))

			for _, mediaType := range virtualMedia.MediaTypes {
				if mediaType == "DVD" {
					log.Info("iLO5 DVD virtual media detected")
					log.Debugf("Virtual DVD: %#v", virtualMedia)
					virtualDevices.Dvd = IloVirtualMedia{
						Id:         virtualMedia.Id,
						MediaTypes: virtualMedia.MediaTypes,
						Actions: IloVirtualMediaActions{
							Insert: virtualMedia.Actions.Insert,
							Eject:  virtualMedia.Actions.Eject,
						},
					}
				}
				if mediaType == "Floppy" {
					log.Info("iLO5 Floppy virtual media detected")
					log.Debugf("Virtual Floppy: %#v", virtualMedia)
					virtualDevices.Floppy = IloVirtualMedia{
						Id:         virtualMedia.Id,
						MediaTypes: virtualMedia.MediaTypes,
						Actions: IloVirtualMediaActions{
							Insert: virtualMedia.Actions.Insert,
							Eject:  virtualMedia.Actions.Eject,
						},
					}
				}
			}
			ilo.VirtualDevices = virtualDevices

		default:
			err := fmt.Errorf("%s is not supported", iloModel)
			log.Error(Wrap(err))
			return err
		}
	}
	log.Debugf("Virtual Devices: %#v", ilo.VirtualDevices)

	return nil
}

// Insert virtual media to HPE iLO
func (ilo *IloClient) InsertVirtualMedia(imageUrl, deviceType string) error {
	log.Info("Insert image into iLO virtual media.")
	// Get ilo virtual mount info
	err := ilo.GetVirtualMedia()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	// Create RedFish client
	c, err := ilo.createRedfishClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	defer c.Logout()

	// Insert target virtual media URL
	var tagertDevice IloVirtualMedia
	if strings.ToLower(deviceType) == "dvd" {
		log.Infof("Insert %s into virtual DVD device.", imageUrl)
		tagertDevice = ilo.VirtualDevices.Dvd
	} else if strings.ToLower(deviceType) == "floppy" {
		log.Infof("Insert %s into virtual Floppy device.", imageUrl)
		tagertDevice = ilo.VirtualDevices.Floppy
	} else {
		err := fmt.Errorf("Unknown device type: %s", deviceType)
		log.Error(Wrap(err))
		return err
	}
	log.Debugf("Target virtual media action endpoint is %v", tagertDevice.Actions.Insert.Target)
	req := IloInsertVirtualMediaReqBody{
		Image: imageUrl,
	}
	res, err := c.Post(tagertDevice.Actions.Insert.Target, req)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	if res.StatusCode >= 400 {
		err := fmt.Errorf("iLO virtual mount failed: %v", res.StatusCode)
		log.Error(Wrap(err))
		return err
	}

	return nil
}

// Eject virtual media on HPE iLO
func (ilo *IloClient) EjectVirtualMedia(deviceType string) error {
	log.Info("Eject image into iLO virtual media.")
	err := ilo.GetVirtualMedia()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}

	// Create RedFish client
	c, err := ilo.createRedfishClient()
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	defer c.Logout()

	// Eject target virtual media URL
	var tagertDevice IloVirtualMedia
	if strings.ToLower(deviceType) == "dvd" {
		log.Info("Eject image from virtual DVD device.")
		tagertDevice = ilo.VirtualDevices.Dvd
	} else if strings.ToLower(deviceType) == "floppy" {
		log.Info("Eject image from virtual Floppy device.")
		tagertDevice = ilo.VirtualDevices.Floppy
	} else {
		err := fmt.Errorf("Unknown device type: %s", deviceType)
		log.Error(Wrap(err))
		return err
	}

	log.Debug(tagertDevice.Actions.Eject.Target)
	var req map[string]string // iLO4 need null body? or json header.
	res, err := c.Post(tagertDevice.Actions.Eject.Target, req)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	if res.StatusCode >= 400 {
		err := fmt.Errorf("Eject virtual media failed: %#v", res)
		log.Error(Wrap(err))
		return err
	}

	return nil
}

func (ilo *IloClient) createRedfishClient() (*gofish.APIClient, error) {
	// Create RedFish client
	config := gofish.ClientConfig{
		Endpoint: "https://" + ilo.Address,
		Session: &gofish.Session{
			Token: ilo.Token,
		},
		TLSHandshakeTimeout: 1,
		Insecure:            true,
		BasicAuth:           false,
	}
	c, err := gofish.Connect(config)
	if err != nil {
		log.Error(Wrap(err))
		return nil, err
	}
	return c, nil
}
