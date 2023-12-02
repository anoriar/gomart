package e2e

import (
	"encoding/json"
	"fmt"
	responses2 "github.com/anoriar/gophermart/internal/gophermart/balance/dto/responses"
	"github.com/anoriar/gophermart/internal/gophermart/order/dto/responses"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
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

	accrualSystemClient *resty.Client
	gophermartClient    *resty.Client
}

func (suite *GophermartSuite) SetupSuite() {
	conf := NewConfig()

	err := env.Parse(conf)

	suite.NoError(err)
	suite.Require().NotEmpty(conf.GophermartBin)
	suite.Require().NotEmpty(conf.GophermartAddress)
	suite.Require().NotEmpty(conf.GophermartRunAddress)
	suite.Require().NotEmpty(conf.DatabaseURI)
	suite.Require().NotEmpty(conf.AccrualSystemBin)
	suite.Require().NotEmpty(conf.AccrualSystemAddress)
	suite.Require().NotEmpty(conf.AccrualSystemRunAddress)

	suite.conf = conf

	suite.startAccrualSystemProcess(*conf)
	suite.startGophermartProcess(*conf)

	suite.accrualSystemClient = resty.New()
	suite.accrualSystemClient.BaseURL = suite.conf.AccrualSystemAddress

	suite.gophermartClient = resty.New()
	suite.gophermartClient.BaseURL = suite.conf.GophermartAddress

}
func (suite *GophermartSuite) startAccrualSystemProcess(conf Config) {
	accrualSystemEnvs := append(os.Environ(),
		"RUN_ADDRESS="+conf.AccrualSystemRunAddress,
	)

	suite.accrualProcess = exec.Command(conf.AccrualSystemBin)
	suite.accrualProcess.Env = append(suite.accrualProcess.Env, accrualSystemEnvs...)

	err := suite.accrualProcess.Start()
	assert.NoError(suite.T(), err, "Error starting the binary")

	time.Sleep(2 * time.Second)
}

func (suite *GophermartSuite) startGophermartProcess(conf Config) {
	gophermartSystemEnvs := append(os.Environ(),
		"RUN_ADDRESS="+conf.GophermartRunAddress,
		"DATABASE_URI="+conf.DatabaseURI,
		"ACCRUAL_SYSTEM_ADDRESS="+conf.AccrualSystemAddress,
	)

	suite.gophermartProcess = exec.Command(conf.GophermartBin)
	suite.gophermartProcess.Env = append(suite.gophermartProcess.Env, gophermartSystemEnvs...)

	stdout, err := os.Create("stdout.log")
	assert.NoError(suite.T(), err)

	defer stdout.Close()
	suite.gophermartProcess.Stdout = stdout

	stderr, err := os.Create("stderr.log")
	assert.NoError(suite.T(), err)

	defer stderr.Close()
	suite.gophermartProcess.Stderr = stderr

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
	suite.Run("ping", func() {

		resp, err := suite.gophermartClient.R().Get("/ping")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Status code:", resp.StatusCode())
		fmt.Println("Response:", resp.String())
	})

	orderNumber := "123456782"

	suite.Run("accrual_add_goods", func() {
		m := []byte(`
			{
			  "match": "Bork",
			  "reward": 10,
			  "reward_type": "%"
			}
		`)

		client := resty.New()
		client.BaseURL = suite.conf.AccrualSystemAddress
		resp, err := suite.accrualSystemClient.R().
			SetHeader("Content-Type", "application/json").
			SetBody(m).
			Post("/api/goods")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusOK, resp.StatusCode())
	})

	suite.Run("accrual_add_order", func() {
		m := []byte(`
				{
				  "order": "` + orderNumber + `",
				  "goods": [
					{
					  "description": "Чайник Bork",
					  "price": 5000
					}
				  ]
				}
		`)

		client := resty.New()
		client.BaseURL = suite.conf.AccrualSystemAddress
		resp, err := suite.accrualSystemClient.R().
			SetHeader("Content-Type", "application/json").
			SetBody(m).
			Post("/api/orders")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusAccepted, resp.StatusCode())
	})

	token := ""
	suite.Run("gophermart_register_user", func() {
		m := []byte(`
				{
				  "login": "test2@gmail.com",
				  "password": "fh49t3jrojf"
				}
		`)

		resp, err := suite.gophermartClient.R().
			SetHeader("Content-Type", "application/json").
			SetBody(m).
			Post("/api/user/register")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusOK, resp.StatusCode())
		token = resp.Header().Get("Authorization")
		suite.NotEmpty(token)
	})

	//удаление пользователя по токену
	defer func() {

		resp, err := suite.gophermartClient.R().
			SetHeader("Authorization", token).
			Delete("/api/user")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusOK, resp.StatusCode())
	}()

	//загрузка заказа
	suite.Run("gophermart_load_order", func() {
		m := []byte(orderNumber)

		resp, err := suite.gophermartClient.R().
			SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", token).
			SetBody(m).
			Post("/api/user/orders")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusAccepted, resp.StatusCode())
	})

	//ожидание, когда подтвердится заказ
	time.Sleep(3 * time.Second)

	suite.Run("gophermart_get_user_orders", func() {

		resp, err := suite.gophermartClient.R().
			SetHeader("Authorization", token).
			Get("/api/user/orders")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusOK, resp.StatusCode())

		var ordersResponse []responses.OrderResponseDto

		err = json.Unmarshal(resp.Body(), &ordersResponse)
		suite.Assert().NoError(err)

		suite.Assert().Equal(1, len(ordersResponse))

		suite.Assert().Equal(orderNumber, ordersResponse[0].Number)
		suite.Assert().Equal("PROCESSED", ordersResponse[0].Status)
		suite.Assert().Equal(500.00, ordersResponse[0].Accrual)

	})

	suite.Run("gophermart_get_user_balance", func() {

		resp, err := suite.gophermartClient.R().
			SetHeader("Authorization", token).
			Get("/api/user/balance")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusOK, resp.StatusCode())

		var balanceResponse responses2.BalanceResponseDto

		err = json.Unmarshal(resp.Body(), &balanceResponse)
		suite.Assert().NoError(err)

		suite.Assert().Equal(500.00, balanceResponse.Current)
		suite.Assert().Equal(0.00, balanceResponse.Withdrawn)

	})

	suite.Run("gophermart_withdraw", func() {
		m := []byte(`
			{
			  "order": "2377225624",
			  "sum": 100
			}
		`)

		resp, err := suite.gophermartClient.R().
			SetHeader("Content-Type", "application/jsons").
			SetHeader("Authorization", token).
			SetBody(m).
			Post("/api/user/balance/withdraw")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusOK, resp.StatusCode())
	})

	suite.Run("gophermart_recheck_user_balance", func() {

		resp, err := suite.gophermartClient.R().
			SetHeader("Authorization", token).
			Get("/api/user/balance")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusOK, resp.StatusCode())

		var balanceResponse responses2.BalanceResponseDto

		err = json.Unmarshal(resp.Body(), &balanceResponse)
		suite.Assert().NoError(err)

		suite.Assert().Equal(400.00, balanceResponse.Current)
		suite.Assert().Equal(100.00, balanceResponse.Withdrawn)

	})

	suite.Run("gophermart_get_user_withdrawals", func() {

		resp, err := suite.gophermartClient.R().
			SetHeader("Authorization", token).
			Get("/api/user/withdrawals")

		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusOK, resp.StatusCode())

		var withdrawalsResponse []responses2.WithdrawalResponseDto

		err = json.Unmarshal(resp.Body(), &withdrawalsResponse)
		suite.Assert().NoError(err)

		suite.Assert().Equal(1, len(withdrawalsResponse))

		suite.Assert().Equal("2377225624", withdrawalsResponse[0].Order)
		suite.Assert().Equal(100.00, withdrawalsResponse[0].Sum)

	})
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(GophermartSuite))
}
