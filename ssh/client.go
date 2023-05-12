package ssh

import (
	"fmt"
	"github.com/a-dakani/LogSpy/configs"
	"github.com/a-dakani/LogSpy/logger"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"os"
	"path/filepath"
)

type Client struct {
	Password string
	Service  configs.Service
	Client   *ssh.Client
	Session  *ssh.Session
}

func (client *Client) DialClient() error {

	var err error
	var privateKey []byte
	var hostWithPort string
	var signer ssh.Signer
	var hostKeyCallback ssh.HostKeyCallback
	conf := &ssh.ClientConfig{}

	if client.Service.Port == 0 {
		client.Service.Port = 22
	}

	hostWithPort = fmt.Sprintf("%s:%d", client.Service.Host, client.Service.Port)

	// TODO check if private key exists
	privateKey, err = os.ReadFile(client.Service.PrivateKeyPath)

	signer, err = ssh.ParsePrivateKey(privateKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	// TODO get KnownHostsPath from config
	// TODO check if known_hosts exists
	hostKeyCallback, err = knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh/known_hosts"))
	if err != nil {
		fmt.Println(err.Error())
	}

	conf = &ssh.ClientConfig{
		User:            client.Service.User,
		HostKeyCallback: hostKeyCallback,
		Auth: []ssh.AuthMethod{
			ssh.Password(client.Password),
			ssh.PublicKeys(signer),
		},
	}

	client.Client, err = ssh.Dial("tcp", hostWithPort, conf)
	if err != nil {
		logger.Fatal(fmt.Sprintf("[%s] unable to connect: %s", client.Service.Host, err))
	}
	defer client.Client.Close()

	return nil
}

func (client *Client) GetSession() error {
	session, err := client.Client.NewSession()
	if err != nil {
		panic(fmt.Sprintf("[%s] unable to create session: %s", client.Service.Host, err))
	}
	client.Session = session
	return err
}
