package ssh

import (
	"bufio"
	"fmt"
	"github.com/a-dakani/remoteLog/configs"
	"github.com/a-dakani/remoteLog/pkg/utils"
	"github.com/jcmturner/gokrb5/v8/config"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"sync"
)

type Spy struct {
	Service  configs.Service
	Client   *ssh.Client
	Sessions []*ssh.Session
}

func (spy *Spy) CreateClient() error {
	var err error

	//hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh/known_hosts"))
	//if err != nil {
	//	logger.Fatal(err.Error())
	//}

	conf := &ssh.ClientConfig{
		User: spy.Service.User,
		//TODO replace InsecureIgnoreHostKey with hostKeyCallback for production
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{},
	}

	//If private key path is provided, use it
	if spy.Service.PrivateKeyPath != "" {
		pemBytes, err := os.ReadFile(spy.Service.PrivateKeyPath)
		if err != nil {
			utils.Warning("Private key file does not exist")
			return err
		}
		// create signer
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			utils.Warning("Private key file is not valid")
			return err
		}
		conf.Auth = append(conf.Auth, ssh.PublicKeys(signer))
		utils.Info(fmt.Sprintf("[%s] Using private key %s ", spy.Service.Host, spy.Service.PrivateKeyPath))
	}
	//If krb5 conf path is provided, use it
	if spy.Service.Krb5ConfPath != "" {
		c, _ := config.Load(spy.Service.Krb5ConfPath)
		//FIXME Error handling is shit for wrong password or unreachable auth server

		sshGSSAPIClient, err := NewKrb5InitiatorClient(spy.Service.User, c)
		if err != nil {
			utils.Warning("Unable to create sshGSSAPIClient")
			return err
		}
		conf.Auth = append(conf.Auth, ssh.GSSAPIWithMICAuthMethod(&sshGSSAPIClient, spy.Service.Host))
		utils.Info(fmt.Sprintf(" [%s] Using krb5 conf %s", spy.Service.Host, spy.Service.Krb5ConfPath))
	}

	spy.Client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", spy.Service.Host, spy.Service.Port), conf)
	if err != nil {
		utils.Warning(fmt.Sprintf("Unable to connect to [%s]", spy.Service.Host))
		return err
	}

	return nil
}

func (spy *Spy) TailFiles() error {
	var wg sync.WaitGroup

	for index, file := range spy.Service.Files {
		wg.Add(1)
		sess, err := spy.Client.NewSession()
		if err != nil {
			return err
		}
		spy.Sessions = append(spy.Sessions, sess)

		sessStdOut, err := sess.StdoutPipe()
		if err != nil {
			return err
		}
		go func() {
			defer wg.Done()
			formattedCopy(os.Stdout, sessStdOut, fmt.Sprintf("[%s]=> ", spy.Service.Files[index].Alias), utils.Red+index)
		}()

		sessStderr, err := sess.StderrPipe()
		if err != nil {
			return err
		}
		go func() {
			defer wg.Done()
			formattedCopy(os.Stderr, sessStderr, fmt.Sprintf("[%s]=> ", spy.Service.Files[index].Alias), utils.Red+index)
		}()

		utils.Info(fmt.Sprintf("[%s] Tailing %s", spy.Service.Host, file.Path))
		go func(index int, path string) {
			defer wg.Done()
			err := spy.Sessions[index].Run(fmt.Sprintf("tail -f %s", path))
			if err != nil {
				utils.Warning(fmt.Sprintf("[%s] Unable to tail %s", spy.Service.Host, path))
				return
			}
			err = spy.Sessions[index].Wait()
			if err != nil {
				utils.Warning(fmt.Sprintf("[%s] Unable to wait for %s", spy.Service.Host, path))
				return
			}
		}(index, file.Path)
	}
	wg.Wait()
	return nil

}

func formattedCopy(dst io.Writer, src io.Reader, appendStr string, color int) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		lineWithAppend := utils.Colorize(appendStr, color) + line + "\n"
		io.WriteString(dst, lineWithAppend)
	}
}

func (spy *Spy) CloseSessions() {
	utils.Info(fmt.Sprintf("[%s] Closing sessions", spy.Service.Host))
	for _, sess := range spy.Sessions {
		sess.Close()
	}
}

func (spy *Spy) CloseClient() {
	utils.Info(fmt.Sprintf("[%s] Closing client", spy.Service.Host))
	spy.Client.Close()

}
