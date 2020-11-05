package slossh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/trenton42/slossh/pkg/recorders"
	"github.com/trenton42/slossh/pkg/session"
	"golang.org/x/crypto/ssh"
)

// Slossh holds the main object
type Slossh struct {
	log        zerolog.Logger
	keyPath    string
	hostKey    ssh.Signer
	recordChan chan session.SlosshSession
	recorders  []recorders.Recorder
}

// New creates a new instance of Slossh
func New() (*Slossh, error) {
	s := Slossh{
		keyPath: "id_rsa",
	}
	s.recorders = recorders.Recorders()
	s.recordChan = make(chan session.SlosshSession, 100)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	s.log = zerolog.New(output).With().Timestamp().Logger()
	err := s.setupConfig()
	return &s, err
}

func (s *Slossh) setupConfig() error {
	var key *rsa.PrivateKey
	var err error
	if s.keyPath != "" {
		data, err := ioutil.ReadFile(s.keyPath)
		if err == nil {
			block, _ := pem.Decode(data)
			var key interface{}
			key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err == nil {
				s.hostKey, err = ssh.NewSignerFromKey(key)
				if err == nil {
					return nil
				}
			}
		}
		s.log.Err(err).Str("key_path", s.keyPath).Msg("Couldn't get stored private key")
	}
	key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return err
	}
	s.hostKey = signer
	if s.keyPath != "" {
		data, err := x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			s.log.Err(err).Msg("Could not marshal the private key")
			return nil
		}
		block := pem.Block{
			Type:  "OPENSSH PRIVATE KEY",
			Bytes: data,
		}
		fp, err := os.OpenFile(s.keyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			s.log.Err(err).Msg("Could open private key for writing")
		}
		err = pem.Encode(fp, &block)
		if err != nil {
			s.log.Err(err).Msg("Could not write out private key")
		}
	}
	return nil
}

// Serve starts the server and waits for connections
func (s *Slossh) Serve(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	go s.Recorder()
	s.log.Info().Int("port", port).Msg("Server started")
	for {
		connection, err := listener.Accept()
		if err != nil {
			s.log.Err(err).Msg("Failed to accept connection")
			continue
		}
		go s.HandleConnection(connection)
	}
}

// HandleConnection does the work on each incoming connection in a new goroutine
func (s *Slossh) HandleConnection(con net.Conn) {
	var remoteIP net.IP
	if addr, ok := con.RemoteAddr().(*net.TCPAddr); ok {
		remoteIP = addr.IP
	}
	s.log.Info().IPAddr("clientIP", remoteIP).Msg("New connection attempt")
	sess := session.NewSession(s.log, s.hostKey)
	sshCon, chans, reqs, err := ssh.NewServerConn(con, sess.Config())
	if err != nil {
		if _, ok := err.(*ssh.ServerAuthError); !ok {
			// only log if some non-authentication error happened. We will always reject all authentication attempts.
			s.log.Err(err).Msg("Handshake failed")
		}
		sess.Close()
		s.recordChan <- *sess
		return
	}
	go ssh.DiscardRequests(reqs)
	for newChannel := range chans {
		s.log.Info().Str("channelType", newChannel.ChannelType()).Msg("Channel attempted")
		newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
	}
	sshCon.Close()
	sess.Close()
	s.recordChan <- *sess
}
