package main

import (
	"flag"
	"github.com/a-dakani/remoteLog/configs"
	"github.com/a-dakani/remoteLog/pkg/ssh"
	"github.com/a-dakani/remoteLog/pkg/utils"
	"reflect"
	"strings"
)

var (
	cfg     configs.Config
	srvs    configs.Services
	srv     configs.Service
	filters []string
)

var service = flag.String("srv", "", "predefined service name in config.services.yaml -srv=myService. This disables the use of -fs, -h, -u, -p, -pk")
var files = flag.String("fs", "", "paths to files being tailed seperated with comma -fs=/var/log/../log-dev-1.log,/var/log/../log-dev-2")
var host = flag.String("h", "", "host to connect to -h=192.168.1.1")
var user = flag.String("u", "", "user to connect to host -u=admin")
var port = flag.Int("p", 22, "port to connect to host -p=22")
var privateKey = flag.String("pk", "", "private key location to connect to host -pk=/home/user/.ssh/id_rsa")
var krb5Conf = flag.String("krb5", "", "krb5.conf location to connect to host -krb5=/etc/krb5.conf")
var filterWords = flag.String("f", "", "filter for the log files -f=ERROR,WARN,FATAL,EXCEPTION")

func init() {
	flag.Parse()

	err := utils.LoadConfig(&cfg)
	if err != nil {
		panic(err)
	}

	if *service != "" {
		err = utils.LoadServices(&srvs)
		if err != nil {
			panic(err)
		}

		for _, confSrv := range srvs.Services {
			if confSrv.Name == *service {
				srv = confSrv
				break
			}
		}

		if reflect.DeepEqual(srv, configs.Service{}) {
			utils.Fatal("Service not found in config.services.yaml")
			panic("Service not found in config.services.yaml")
		}

	} else {
		srv = configs.Service{
			Name:           "ArgService",
			Host:           *host,
			User:           *user,
			Port:           *port,
			PrivateKeyPath: *privateKey,
			Krb5ConfPath:   *krb5Conf,
			Files:          configs.ParseFiles(*files),
		}
		if _, err = srv.IsFullyConfigured(); err != nil {
			utils.ProcessArgumentError()
			panic(err)
		}
	}
	filters = strings.Split(*filterWords, ",")

	if filters[0] == "" && len(filters) == 1 {
		utils.Warning("No filter words provided. Proceeding without filters")
	} else {
		utils.Info("Filter words provided:" + *filterWords)
	}
}

func main() {
	s := ssh.Spy{
		Service: srv,
	}
	err := s.CreateClient()
	if err != nil {
		utils.Fatal(err.Error())
		panic(err)
	}
	defer s.CloseClient()

	err = s.TailFiles()
	if err != nil {
		utils.Warning(err.Error())
		panic(err)
	}
	defer s.CloseSessions()

}
