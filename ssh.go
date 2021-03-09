/*
------------------------------------------------------------------------------------------------------------------------
####### uuid ####### (c) 2020-2021 mls-361 ######################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package ssh

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/mls-361/logger"
	"golang.org/x/crypto/ssh"
)

const (
	_defaultTimeout = 5 * time.Second
)

type (
	// ClientOptions AFAIRE.
	ClientOptions struct {
		Host       string
		Port       int
		Username   string
		Password   string
		KeyFile    string
		Passphrase string
		Timeout    time.Duration
	}
)

// NewClient AFAIRE.
func (co *ClientOptions) NewClient() *Client {
	if co.Port == 0 {
		co.Port = 22
	}

	if co.Timeout == 0 {
		co.Timeout = _defaultTimeout
	}

	return &Client{options: co}
}

type (
	// Client AFAIRE.
	Client struct {
		options *ClientOptions
	}
)

func (c *Client) readKeyFile() (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(c.options.KeyFile)
	if err != nil {
		return nil, err
	}

	if c.options.Passphrase != "" {
		key, err := ssh.ParseRawPrivateKeyWithPassphrase(buf, []byte(c.options.Passphrase))
		if err != nil {
			return nil, err
		}

		return ssh.NewSignerFromKey(key)
	}

	return ssh.ParsePrivateKey(buf)
}

func (c *Client) configure() (*ssh.ClientConfig, error) {
	auths := []ssh.AuthMethod{}

	if c.options.Password != "" {
		auths = append(auths, ssh.Password(c.options.Password))
	}

	if c.options.KeyFile != "" {
		key, err := c.readKeyFile()
		if err != nil {
			return nil, err
		}

		auths = append(auths, ssh.PublicKeys(key))
	}

	cfg := &ssh.ClientConfig{
		User:            c.options.Username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         c.options.Timeout,
	}

	return cfg, nil
}

// Connect AFAIRE.
func (c *Client) Connect(logger logger.Logger) (*Connection, error) {
	cfg, err := c.configure()
	if err != nil {
		return nil, err
	}

	ssh, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.options.Host, c.options.Port), cfg)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		logger: logger,
		client: c,
		ssh:    ssh,
	}

	return conn, nil
}

type (
	// Connection AFAIRE.
	Connection struct {
		logger logger.Logger
		client *Client
		ssh    *ssh.Client
	}
)

// Host AFAIRE.
func (c *Connection) Host() string {
	return c.client.options.Host
}

// Port AFAIRE.
func (c *Connection) Port() int {
	return c.client.options.Port
}

// Username AFAIRE.
func (c *Connection) Username() string {
	return c.client.options.Username
}

// NewSession AFAIRE.
func (c *Connection) NewSession() (*Session, error) {
	s, err := c.ssh.NewSession()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Session: s,
		logger:  c.logger,
		client:  c.client,
	}

	return session, nil
}

// ReadStream AFAIRE.
func (c *Connection) ReadStream(cmd string, timeout time.Duration) (*Stream, error) {
	session, err := c.NewSession()
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
		return nil, err
	}

	stderrReader, err := session.StderrPipe()
	if err != nil {
		return nil, err
	}

	stderrScanner := bufio.NewScanner(io.MultiReader(stderrReader))
	stdoutScanner := bufio.NewScanner(io.MultiReader(stdoutReader))

	if err := session.Start(cmd); err != nil {
		return nil, err
	}

	go s.readData(timeout, stderrScanner, stdoutScanner) //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

	return s, nil
}

// Disconnect AFAIRE.
func (c *Connection) Disconnect() {
	c.ssh.Close()
}

type (
	// Session AFAIRE.
	Session struct {
		*ssh.Session
		logger logger.Logger
		client *Client
	}
)

func (s *Session) trace(cmd string) {
	if s.logger != nil {
		s.logger.Debug( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"SSH",
			"server", s.client.options.Host,
			"username", s.client.options.Username,
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

type (
	// Stream AFAIRE.
	Stream struct {
		session *Session
		stderr  chan string
		stdout  chan string
		done    chan bool
		err     error
	}
)

func (s *Stream) readData(timeout time.Duration, stderrScanner, stdoutScanner *bufio.Scanner) {
	defer close(s.done)

	defer close(s.stdout)
	defer close(s.stderr)

	defer s.session.Close()

	stop := make(chan struct{}, 1)
	group := sync.WaitGroup{}

	group.Add(2)

	go func() { //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		for stderrScanner.Scan() {
			s.stderr <- stderrScanner.Text()
		}

		group.Done()
	}()

	go func() { //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		for stdoutScanner.Scan() {
			s.stdout <- stdoutScanner.Text()
		}

		group.Done()
	}()

	go func() { //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		group.Wait()
		stop <- struct{}{}
	}()

	select {
	case <-stop:
		s.err = s.session.Wait()
		s.done <- true
	case <-time.After(timeout):
		s.done <- false
	}
}

// Stderr AFAIRE.
func (s *Stream) Stderr() <-chan string {
	return s.stderr
}

// Stdout AFAIRE.
func (s *Stream) Stdout() <-chan string {
	return s.stdout
}

// Done AFAIRE.
func (s *Stream) Done() <-chan bool {
	return s.done
}

// Err AFAIRE.
func (s *Stream) Err() error {
	return s.err
}

/*
######################################################################################################## @(°_°)@ #######
*/
