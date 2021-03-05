/*
------------------------------------------------------------------------------------------------------------------------
####### uuid ####### (c) 2020-2021 mls-361 ######################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package ssh

import (
	"github.com/mls-361/failure"
	"github.com/mls-361/logger"
)

type (
	// Clients AFAIRE.
	Clients map[string]map[string]*Client
)

// NewClients AFAIRE.
func NewClients(cos []*ClientOptions) Clients {
	clients := Clients{}

	for _, co := range cos {
		_, ok := clients[co.Host]
		if !ok {
			clients[co.Host] = make(map[string]*Client)
		}

		clients[co.Host][co.Username] = co.NewClient()
	}

	return clients
}

// Connect AFAIRE.
func (c Clients) Connect(host, username string, logger *logger.Logger) (*Connection, error) {
	client, ok := c[host][username]
	if !ok {
		return nil,
			failure.New(nil).
				Set("server", host).
				Set("user", username).
				Msg("this SSH server or user does not exist") //////////////////////////////////////////////////////////
	}

	return client.Connect(logger)
}

/*
######################################################################################################## @(°_°)@ #######
*/
