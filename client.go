/*
------------------------------------------------------------------------------------------------------------------------
####### ssh ####### (c) 2020-2021 mls-361 ########################################################## MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package ssh

import (
	"github.com/mls-361/logger"
	"golang.org/x/crypto/ssh"
)

type (
	// Client AFAIRE.
	Client struct {
		host     string
		port     int
		username string
		addr     string
		config   *ssh.ClientConfig
	}
)

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
