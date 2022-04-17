package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	CERT_PASSWORD        = 1
	CERT_PUBLIC_KEY_FILE = 2
	DEFAULT_TIMEOUT      = 3 // seconds
)

type SSHClient struct {
	Ip   string
	User string
	Cert string // password or key file path
	Port int
}

type Connection struct {
	session *ssh.Session
	client  *ssh.Client
}

func (ssh_client *SSHClient) publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("PublicKey read error:", err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatal("PrivateKey parse error:", err)
	}
	return ssh.PublicKeys(key)
}

func (ssh_client *SSHClient) Connect(mode int) Connection {

	var ssh_config *ssh.ClientConfig
	var auth []ssh.AuthMethod

	switch mode {
	case CERT_PASSWORD:
		auth = []ssh.AuthMethod{ssh.Password(ssh_client.Cert)}
	case CERT_PUBLIC_KEY_FILE:
		auth = []ssh.AuthMethod{ssh_client.publicKeyFile(ssh_client.Cert)}
	default:
		log.Fatal("Mode is not supported: ", mode)
	}

	ssh_config = &ssh.ClientConfig{
		User: ssh_client.User,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * DEFAULT_TIMEOUT,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ssh_client.Ip, ssh_client.Port), ssh_config)
	if err != nil {
		log.Fatal(err)
	}

	connection := Connection{
		session: nil,
		client:  client,
	}

	connection.newSession()

	return connection
}

func (connection *Connection) RunCmd(cmd string) (error, string) {
	connection.newSession()
	out, err := connection.session.CombinedOutput(cmd)
	return err, string(out)
}

func (connection *Connection) Close() {
	if connection.session != nil {
		connection.session.Close()
	}
	connection.client.Close()
}

func (connection *Connection) newSession() {
	if connection.session != nil {
		connection.session.Close()
	}
	var err error
	connection.session, err = connection.client.NewSession()
	if err != nil {
		connection.Close()
		log.Fatal("Can't create new session: ", err)
	}
}
