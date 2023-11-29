package e2e

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"testing"
	"time"
)

type GophermartSuite struct {
	suite.Suite

	gophermartProcess *exec.Cmd

	accrualProcess *exec.Cmd
	conf           *Config
}

func (suite *GophermartSuite) SetupSuite() {
	conf := NewConfig()

	err := env.Parse(conf)

	suite.NoError(err)
	suite.Require().NotEmpty(conf.GophermartBin)
	suite.Require().NotEmpty(conf.GophermartAddress)
	suite.Require().NotEmpty(conf.DatabaseURI)
	suite.Require().NotEmpty(conf.AccrualSystemBin)
	suite.Require().NotEmpty(conf.AccrualSystemAddress)

	suite.conf = conf

	accrualSystemEnvs := append(os.Environ(),
		"RUN_ADDRESS=localhost:8081",
	)

	suite.accrualProcess = exec.Command(conf.AccrualSystemBin)
	suite.accrualProcess.Env = append(suite.accrualProcess.Env, accrualSystemEnvs...)

	err = suite.accrualProcess.Start()
	assert.NoError(suite.T(), err, "Error starting the binary")

	time.Sleep(2 * time.Second)

	gophermartSystemEnvs := append(os.Environ(),
		"RUN_ADDRESS=localhost:8080",
		"DATABASE_URI="+conf.DatabaseURI,
		"ACCRUAL_SYSTEM_ADDRESS="+conf.AccrualSystemAddress,
	)

	suite.gophermartProcess = exec.Command(conf.GophermartBin)
	suite.gophermartProcess.Env = append(suite.accrualProcess.Env, gophermartSystemEnvs...)

	err = suite.gophermartProcess.Start()
	assert.NoError(suite.T(), err, "Error starting the binary")

	time.Sleep(2 * time.Second)
}

func (suite *GophermartSuite) TearDownSuite() {
	if suite.accrualProcess != nil && suite.accrualProcess.Process != nil {
		_ = suite.gophermartProcess.Process.Kill() // Kill the process when the suite is done
	}

	if suite.gophermartProcess != nil && suite.gophermartProcess.Process != nil {
		_ = suite.gophermartProcess.Process.Kill() // Kill the process when the suite is done
	}
}

func (suite *GophermartSuite) TestGophermart() {
	suite.Run("test", func() {
		client := resty.New()
		client.BaseURL = suite.conf.GophermartAddress

		resp, err := client.R().Get("/ping")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Status code:", resp.StatusCode())
		fmt.Println("Response:", resp.String())
	})
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(GophermartSuite))
}
