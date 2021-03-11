/*
------------------------------------------------------------------------------------------------------------------------
####### ssh ####### (c) 2020-2021 mls-361 ########################################################## MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package ssh

import (
	"github.com/mls-361/failure"
	"github.com/mls-361/logger"
)

type (
	// Crypto AFAIRE.
	Crypto interface {
		DecryptString(text string) (string, error)
	}

	// Clients AFAIRE.
	Clients map[string]map[string]*Client
)

// NewClients AFAIRE.
func NewClients(aofCfg []*Config, crypto Crypto) (Clients, error) {
	clients := Clients{}

	for _, cfg := range aofCfg {
		if crypto != nil {
			if cfg.Password != "" {
				if p, err := crypto.DecryptString(cfg.Password); err != nil {
					return nil, err
				} else {
					cfg.Password = p
				}
			}

			if cfg.Passphrase != "" {
				if p, err := crypto.DecryptString(cfg.Passphrase); err != nil {
					return nil, err
				} else {
					cfg.Passphrase = p
				}
			}
		}

		_, ok := clients[cfg.Host]
		if !ok {
			clients[cfg.Host] = make(map[string]*Client)
		}

		client, err := cfg.NewClient()
		if err != nil {
			return nil, err
		}

		clients[cfg.Host][cfg.Username] = client
	}

	return clients, nil
}

// Connect AFAIRE.
func (c Clients) Connect(host, username string, logger logger.Logger) (*Connection, error) {
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
