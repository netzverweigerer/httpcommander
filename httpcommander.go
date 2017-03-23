package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type CommandProperties struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type Configuration struct {
	ListenAddress string                       `json:"listenAddress"`
	TlsKey        string                       `json:"tlsKey"`
	TlsCert       string                       `json:"tlsCert"`
	CaCert        string                       `json:"caCert"`
	CommandMap    map[string]CommandProperties `json:"commands"`
}

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		log.Println("missing comandline arg for configfile.")
		os.Exit(1)
	}

	log.Println("using config: " + args[0])
	configfile, _ := os.Open(args[0])
	decoder := json.NewDecoder(configfile)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)

	}

	var usetls = false
	var useclientauth = false
	var tlskey string = ""
	var tlscert string = ""
	//var cacert = ""

	if len(configuration.TlsCert) > 0 && len(configuration.TlsKey) > 0 {
		usetls = true
		tlskey = configuration.TlsKey
		tlscert = configuration.TlsCert
	} else {
		usetls = false
	}

	if len(configuration.CaCert) > 0 {
		useclientauth = true
	} else {
		useclientauth = false
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/cmd/:cmd", func(c *gin.Context) {
		cmdstr := c.Param("cmd")
		if _, ok := configuration.CommandMap[cmdstr]; ok {
			log.Println(configuration.CommandMap[cmdstr].Args)
			cmdargs := configuration.CommandMap[cmdstr].Args
			cmdstr := configuration.CommandMap[cmdstr].Command

			log.Printf("Found cmd for : %s -> %s", cmdstr, cmdargs)
			cmd := exec.Command(cmdstr, cmdargs...)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Println(err)
				c.JSON(500, gin.H{
					"cmd":    c.Param("cmd"),
					"output": err.Error(),
				})
			} else {
				c.JSON(200, gin.H{
					"cmd":    c.Param("cmd"),
					"output": out.String(),
				})

			}
		} else {
			c.JSON(404, gin.H{
				"cmd":    c.Param("cmd"),
				"output": "Command not found in map.",
			})
		}

		log.Println("done.")
	})

	if useclientauth {
		// Load CA cert
		caCert, err := ioutil.ReadFile(configuration.CaCert)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		server := &http.Server{
			Addr:    configuration.ListenAddress,
			Handler: router,
			TLSConfig: &tls.Config{
				ClientAuth: tls.RequireAndVerifyClientCert,
				ClientCAs:  caCertPool,
			},
		}
		log.Printf("Starting tls server with auth on %s\n", configuration.ListenAddress)
		err = server.ListenAndServeTLS(
			tlscert,
			tlskey)
		if err != nil {
			log.Fatal(err)
		}
	} else {

		if usetls {
			log.Printf("Starting tls server. Listen on %s\n", configuration.ListenAddress)
			router.RunTLS(configuration.ListenAddress, tlscert, tlskey)
		} else {
			log.Printf("Starting server. Listen on %s\n", configuration.ListenAddress)
			router.Run(configuration.ListenAddress)
		}
	}
}
