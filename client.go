/*
------------------------------------------------------------------------------------------------------------------------
####### ssh ####### (c) 2020-2021 mls-361 ########################################################## MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package ssh

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mls-361/failure"
	"github.com/mls-361/logger"
	"golang.org/x/crypto/ssh"
)

type (
	// Config AFAIRE.
	Config struct {
		Host       string
		Port       int
		Username   string
		Password   string
		KeyFile    string
		Passphrase string
		Timeout    time.Duration
	}

	// Client AFAIRE.
	Client struct {
		host     string
		port     int
		username string
		addr     string
		config   *ssh.ClientConfig
	}
)

func (cfg *Config) signer() (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(cfg.KeyFile)
	if err != nil {
		return nil, err
	}

	if cfg.Passphrase != "" {
		key, err := ssh.ParseRawPrivateKeyWithPassphrase(buf, []byte(cfg.Passphrase))
		if err != nil {
			return nil, err
		}

		return ssh.NewSignerFromKey(key)
	}

	return ssh.ParsePrivateKey(buf)
}

func (cfg *Config) configure() (*ssh.ClientConfig, error) {
	auths := []ssh.AuthMethod{}

	if cfg.KeyFile != "" {
		key, err := cfg.signer()
		if err != nil {
			return nil, err
		}

		auths = append(auths, ssh.PublicKeys(key))
	}

	// Nécessite "PasswordAuthentication yes" in /etc/ssh/sshd_config
	if cfg.Password != "" {
		auths = append(auths, ssh.Password(cfg.Password))
	}

	scc := &ssh.ClientConfig{
		User:            cfg.Username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // AFINIR
		Timeout:         cfg.Timeout,
	}

	return scc, nil
}

// NewClient AFAIRE.
func (cfg *Config) NewClient() (*Client, error) {
	if cfg.Host == "" {
		return nil, failure.New(nil).Msg("the Host field cannot be empty") /////////////////////////////////////////////
	}

	if cfg.Username == "" {
		return nil, failure.New(nil).Msg("the Username field cannot be empty") /////////////////////////////////////////
	}

	config, err := cfg.configure()
	if err != nil {
		return nil, err
	}

	port := cfg.Port

	if port == 0 {
		port = 22
	}

	c := &Client{
		host:     cfg.Host,
		port:     port,
		username: cfg.Username,
		addr:     fmt.Sprintf("%s:%d", cfg.Host, port),
		config:   config,
	}

	return c, nil
}

// Connect AFAIRE.
func (cl *Client) Connect(logger logger.Logger) (*Connection, error) {
	ssh, err := ssh.Dial("tcp", cl.addr, cl.config)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		client: cl,
		logger: logger,
		ssh:    ssh,
	}

	return conn, nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
