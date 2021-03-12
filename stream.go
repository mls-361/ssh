/*
------------------------------------------------------------------------------------------------------------------------
####### ssh ####### (c) 2020-2021 mls-361 ########################################################## MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package ssh

import (
	"bufio"
	"sync"
	"time"
)

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

// Close AFAIRE.
func (s *Stream) Close() {
	if s.done != nil {
		close(s.done)
		s.done = nil
	}

	if s.stdout != nil {
		close(s.stdout)
		s.stdout = nil
	}

	if s.stderr != nil {
		close(s.stderr)
		s.stderr = nil
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
