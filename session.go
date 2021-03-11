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
	// Session AFAIRE.
	Session struct {
		*ssh.Session
		client *Client
		logger logger.Logger
	}
)

func (s *Session) trace(cmd string) {
	if s.logger != nil {
		s.logger.Debug( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"SSH",
			"server", s.client.host,
			"username", s.client.username,
			"cmd", cmd,
		)
	}
}

// CombinedOutput AFAIRE.
func (s *Session) CombinedOutput(cmd string) ([]byte, error) {
	s.trace(cmd)
	return s.Session.CombinedOutput(cmd)
}

// Output AFAIRE.
func (s *Session) Output(cmd string) ([]byte, error) {
	s.trace(cmd)
	return s.Session.Output(cmd)
}

// Run AFAIRE.
func (s *Session) Run(cmd string) error {
	s.trace(cmd)
	return s.Session.Run(cmd)
}

// Start AFAIRE.
func (s *Session) Start(cmd string) error {
	s.trace(cmd)
	return s.Session.Start(cmd)
}

/*
######################################################################################################## @(°_°)@ #######
*/
