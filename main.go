package main

import (
	"flag"
	"github.com/a-dakani/LogSpy/configs"
	"github.com/a-dakani/LogSpy/logger"
	"github.com/a-dakani/LogSpy/ssh"
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
var filterWords = flag.String("f", "", "filter for the log files -f=ERROR,WARN,FATAL,EXCEPTION")

func init() {
	flag.Parse()
	configs.LoadConfig(&cfg)
	if *service != "" {
		configs.LoadServices(&srvs)
		for _, confSrv := range srvs.Services {
			if confSrv.Name == *service {
				srv = confSrv
				break
			}
		}
		if reflect.DeepEqual(srv, configs.Service{}) {
			logger.Fatal("Service not found in config.services.yaml")
		}
	} else {
		srv = configs.Service{
			Name:           "ArgService",
			Host:           *host,
			User:           *user,
			Port:           *port,
			PrivateKeyPath: *privateKey,
			Files:          strings.Split(*files, ","),
		}
		if !srv.IsFullyConfigured() {
			logger.ProcessArgumentError()

		}
	}
	filters = strings.Split(*filterWords, ",")

	if filters[0] == "" && len(filters) == 1 {
		logger.Warning("No filter words provided. Proceeding without filters")
	} else {
		logger.Info("Filter words provided:" + *filterWords)
	}
}

func main() {
	client := ssh.Client{
		Service:  srv,
		Password: "test",
	}
	err := client.DialClient()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	err = client.GetSession()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	output, err := client.Session.CombinedOutput("pwd")
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	logger.Info(string(output))
}
