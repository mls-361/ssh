/*
------------------------------------------------------------------------------------------------------------------------
####### ssh ####### (c) 2020-2021 mls-361 ########################################################## MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package ssh

import (
	"bufio"
	"io"
	"time"

	"github.com/mls-361/logger"
	"golang.org/x/crypto/ssh"
)

type (
	// Connection AFAIRE.
	Connection struct {
		client *Client
		logger logger.Logger
		ssh    *ssh.Client
	}
)

// Host AFAIRE.
func (conn *Connection) Host() string {
	return conn.client.host
}

// Port AFAIRE.
func (conn *Connection) Port() int {
	return conn.client.port
}

// Username AFAIRE.
func (conn *Connection) Username() string {
	return conn.client.username
}

// NewSession AFAIRE.
func (conn *Connection) NewSession() (*Session, error) {
	s, err := conn.ssh.NewSession()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Session: s,
		client:  conn.client,
		logger:  conn.logger,
	}

	return session, nil
}

// ReadStream AFAIRE.
func (conn *Connection) ReadStream(cmd string, timeout time.Duration) (*Stream, error) {
	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	s := &Stream{
		session: session,
		stderr:  make(chan string),
		stdout:  make(chan string),
		done:    make(chan bool),
	}

	stdoutReader, err := session.StdoutPipe()
	if err != nil {
		s.Close()
		return nil, err
	}

	stderrReader, err := session.StderrPipe()
	if err != nil {
		s.Close()
		return nil, err
	}

	stderrScanner := bufio.NewScanner(io.MultiReader(stderrReader))
	stdoutScanner := bufio.NewScanner(io.MultiReader(stdoutReader))

	if err := session.Start(cmd); err != nil {
		s.Close()
		return nil, err
	}

	go s.readData(timeout, stderrScanner, stdoutScanner) //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

	return s, nil
}

// Disconnect AFAIRE.
func (conn *Connection) Disconnect() {
	conn.ssh.Close()
}

/*
######################################################################################################## @(°_°)@ #######
*/
