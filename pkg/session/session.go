package session

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

// ID wraps an ssh.SessionID for better ecoding
type ID []byte

// MarshalJSON turns it into a hex string
func (s ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// String satisties the stringer type
func (s ID) String() string {
	return hex.EncodeToString(s)
}

// JSONKey is a wrapper around an ssh key that is more json friendly
type JSONKey struct {
	ssh.PublicKey
}

// MarshalJSON turns the key into something usable
func (j JSONKey) MarshalJSON() ([]byte, error) {
	if j.PublicKey == nil {
		return json.Marshal(nil)
	}
	tmp := struct {
		Key         string
		Fingerprint string
		Type        string
	}{
		Key:         strings.TrimSpace(string(ssh.MarshalAuthorizedKey(j.PublicKey))),
		Fingerprint: ssh.FingerprintSHA256(j.PublicKey),
		Type:        j.Type(),
	}

	return json.Marshal(tmp)
}

// SlosshSession captures information about a session
type SlosshSession struct {
	SessionID     ID
	IP            net.IP
	ClientVersion string
	Attempts      []Attempt
	Start         time.Time
	Finish        time.Time
	log           zerolog.Logger
	publicKey     ssh.Signer
}

// Attempt holds an ssh attempt, either through a public key or through username and password
type Attempt struct {
	Username string
	Key      JSONKey
	Password string
}

// NewSession initializes a session
func NewSession(l zerolog.Logger, key ssh.Signer) *SlosshSession {
	s := SlosshSession{
		Start:     time.Now(),
		log:       l,
		publicKey: key,
		Attempts:  []Attempt{},
	}
	return &s
}

func (s *SlosshSession) passwordCallback(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	s.Attempts = append(s.Attempts, Attempt{Username: c.User(), Password: string(pass)})
	return nil, fmt.Errorf("password rejected")
}
func (s *SlosshSession) publicKeyCallback(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
	s.Attempts = append(s.Attempts, Attempt{Username: c.User(), Key: JSONKey{pubKey}})
	return nil, fmt.Errorf("public key rejected")
}

func (s *SlosshSession) bannerCallback(c ssh.ConnMetadata) string {
	return "*** Authorized access only ***"
}

func (s *SlosshSession) authLogCallback(c ssh.ConnMetadata, method string, err error) {
	s.ClientVersion = string(c.ClientVersion())
	s.SessionID = c.SessionID()
	if ip, ok := c.RemoteAddr().(*net.TCPAddr); ok {
		s.IP = ip.IP
	}
	s.log.Debug().Str("session-id", ID(c.SessionID()).String()).Bytes("client-version", c.ClientVersion()).Int("attempts", len(s.Attempts)).Str("method", method).Msg("Login Attempt")
}

// Config returns an ssh config that will feed information into the session
func (s *SlosshSession) Config() *ssh.ServerConfig {
	out := ssh.ServerConfig{
		ServerVersion:     "SSH-2.0-OpenSSH_8.1",
		PasswordCallback:  s.passwordCallback,
		PublicKeyCallback: s.publicKeyCallback,
		AuthLogCallback:   s.authLogCallback,
	}
	out.AddHostKey(s.publicKey)
	return &out
}

// Close completes the session
func (s *SlosshSession) Close() {
	s.Finish = time.Now()
	s.log.Info().Int("login_attempts", len(s.Attempts)).Dur("duration", s.Finish.Sub(s.Start)).Msg("Session closed")
}
