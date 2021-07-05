package driver

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
)

type Server struct {
	Address       string `yaml:"address"`
	KsBaseUrl     string `yaml:"kickstart-base-url"`
	OsUrl         string `yaml:"os-url"`
	RootPassword  string `default:"password" yaml:"root-password"`
	KsUrl         string
	SshPublicKey  string
	SshPrivateKey string
	Hostname      string
}

const (
	defaultShellTimeout    = 1
	defaultWebTimeout      = 1
	defaultInstallTimeout  = 3600 //sec
	defaultInstallInterval = 30
	defaultTicker          = 15 //sec
)

func (s *Server) Validate() error {
	// Check os image iso URL
	client := &http.Client{
		Timeout: defaultWebTimeout * time.Second,
	}
	imageUrl := s.OsUrl
	respOsImage, err := client.Get(imageUrl)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	if respOsImage.StatusCode >= 400 {
		err := fmt.Errorf("Could not access %s: %d", imageUrl, respOsImage.StatusCode)
		log.Error(Wrap(err))
		return err
	}

	// Check ks image iso URL
	ksUrl := s.KsUrl
	respKsImage, err := client.Get(ksUrl)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	if respKsImage.StatusCode >= 400 {
		err := fmt.Errorf("Could not access %v: %d", s.OsUrl, respKsImage.StatusCode)
		log.Error(Wrap(err))
		return err
	}

	return nil
}

func (s *Server) RemoteShell(shell string, port int) error {
	address := s.Address
	sshClient, err := ssh.NewNativeClient(
		"root",
		address,
		port,
		&ssh.Auth{
			Passwords: []string{
				s.RootPassword,
			},
		},
	)
	log.Debugf("Initialize ssh client: %#v", sshClient)
	if err != nil {
		log.Error(Wrap(err))
		return err
	}
	log.Debugf("Execute %s via ssh", shell)
	sshClient.(*ssh.NativeClient).Config.Timeout = defaultShellTimeout * time.Second

	err = sshClient.Shell(shell)
	if err != nil {
		log.Debug(Wrap(err))
		return err
	}

	return nil
}

func (s *Server) WaitOsInstallation() error {
	log.Infof("Waiting for OS installation")
	log.Infof("Trying to ssh access to new server... Timeout is %v sec", defaultInstallTimeout)
	if err := mcnutils.WaitForSpecific(s.sshAvailableFunc(), defaultInstallTimeout/defaultInstallInterval, defaultInstallInterval*time.Second); err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}

func (s *Server) sshAvailableFunc() func() bool {
	return func() bool {
		log.Infof("Waiting for SSH to be available...")
		if err := s.RemoteShell("echo hello", 22); err != nil {
			return false
		}
		return true
	}
}

func (s *Server) CopySshPubKey() error {
	pubkey := s.SshPublicKey
	log.Debugf("CopyPubKey: %s", pubkey)
	log.Debugf("TargetServer: %s", s.Address)
	if pubkey == "" {
		err := fmt.Errorf("Public key is not generated")
		log.Error(Wrap(err))
		return err
	}
	shell := fmt.Sprintf(`
	mkdir -p /root/.ssh 
	echo "%s docker-machine-ov" >> /root/.ssh/authorized_keys
	chmod 0600 /root/.ssh/authorized_keys
	`, pubkey)
	log.Debugf("Shell: %s", shell)
	if err := s.RemoteShell(shell, 22); err != nil {
		log.Error(Wrap(err))
		return err
	}
	return nil
}
