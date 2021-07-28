package driver

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/docker/machine/libmachine/mcnflag"
	"gopkg.in/yaml.v2"
)

type HpeConfig struct {
	Oneview *Oneview `yaml:"oneview"`
	Server  *Server  `yaml:"server"`
	Yaml    *Yaml
}

type Yaml struct {
	Path string
}

// Read config yaml
func (y *Yaml) Read() (*HpeConfig, error) {
	bytes, err := ReadFile(y.Path)
	if err != nil {
		return nil, err
	}

	var conf HpeConfig
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func ReadFile(path string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func Wrap(err error) error {
	if err == nil {
		return err
	}

	// Get stack trace
	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	s := fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)

	// Wrap and return including stack trace
	return fmt.Errorf("%w %s", err, s)
}

/**************************************
Prameters for HPE Servers
***************************************/
var mcnFlags = []mcnflag.Flag{
	/**************
	HPE OneView setting
	**************/
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_YAML",
		Name:   driverName + "-yaml",
		Usage:  "(Option) Configuration YAML file path",
		Value:  "",
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_ONEVIEW_ENDPOINT",
		Name:   driverName + "-oneview-endpoint",
		Usage:  "HPE OneView endpoint URL.",
		Value:  "",
	},
	mcnflag.IntFlag{
		EnvVar: strings.ToUpper(driverName) + "_ONEVIEW_API_VERSION",
		Name:   driverName + "-oneview-api-version",
		Usage:  "HPE OneView API version.",
		Value:  1800,
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_ONEVIEW_USER",
		Name:   driverName + "-oneview-user",
		Usage:  "HPE OneView user name. The user should be infrastructure administrator.",
		Value:  "administrator",
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_ONEVIEW_PASSWORD",
		Name:   driverName + "-oneview-password",
		Usage:  "HPE OneView user password.",
		Value:  "password",
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_ONEVIEW_DOMAIN",
		Name:   driverName + "-oneview-domain",
		Usage:  "(Option) HPE OneView domain name.",
		Value:  "",
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_ONEVIEW_SERVER_PROFILE_TEMPLATE",
		Name:   driverName + "-oneview-server-profile-template",
		Usage:  "HPE OneView server profile template name used when creating target server.",
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_ONEVIEW_SERVER_HARDWARE",
		Name:   driverName + "-oneview-server-hardware",
		Usage:  "HPE OneView server hardware name. This server will be target server. (EXACTLY same name as OneView displayed. There is a case to need spaces between strings when hardware name is displayed with sapces in OneView.)",
	},
	/**************
	New server setting
	**************/
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_SERVER_ADDRESS",
		Name:   driverName + "-server-address",
		Usage:  "Target server IP address.",
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_SERVER_ROOT_PASSWORD",
		Name:   driverName + "-server-root-password",
		Usage:  "Target server root user password.",
		Value:  "password",
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_SERVER_KICKSTART_BASE_URL",
		Name:   driverName + "-server-kickstart-base-url",
		Usage:  "Kickstart image base URL. If your kickstart iso image is on http://web01/docker/kickstart.iso, you shoud set this value as http://web01/docker.",
	},
	mcnflag.StringFlag{
		EnvVar: strings.ToUpper(driverName) + "_SERVER_OS_URL",
		Name:   driverName + "-server-os-url",
		Usage:  "OS image URL.",
	},
	/**************
	Common
	**************/
	mcnflag.BoolFlag{
		EnvVar: strings.ToUpper(driverName) + "_DEBUG",
		Name:   driverName + "-debug",
		Usage:  "(Option) Debug flag for this driver.",
	},
}
