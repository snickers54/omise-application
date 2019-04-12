package internal

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	omise "github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	log "github.com/sirupsen/logrus"
)

type Positions struct {
	Name   int
	Amount int
	Year   int
	Month  int
	Card   int
}

type WorkerPool struct {
	Wg           sync.WaitGroup
	NbWorker     int
	ChannelWork  chan *[]string
	ChannelsQuit []chan bool
	Stats
	Positions
}

func NewWorkerPool() *WorkerPool {
	nbWorker, err := strconv.Atoi(os.Getenv("NB_WORKER"))
	if err != nil {
		nbWorker = 10
	}
	posAmount, _ := strconv.Atoi(os.Getenv("POS_AMOUNT"))
	posName, _ := strconv.Atoi(os.Getenv("POS_NAME"))
	posCard, _ := strconv.Atoi(os.Getenv("POS_CARD"))
	posYear, _ := strconv.Atoi(os.Getenv("POS_YEAR"))
	posMonth, _ := strconv.Atoi(os.Getenv("POS_MONTH"))
	return &WorkerPool{
		Wg:          sync.WaitGroup{},
		NbWorker:    nbWorker,
		ChannelWork: make(chan *[]string),
		Positions: Positions{
			Amount: posAmount,
			Name:   posName,
			Card:   posCard,
			Month:  posMonth,
			Year:   posYear,
		},
		Stats: Stats{
			NbDonations:  0,
			NbSuccessful: 0,
			TotalAmount:  big.NewInt(0),
			TotalFaulty:  big.NewInt(0),
		},
	}
}

func (c *WorkerPool) Close() {
	close(c.ChannelWork)
}

func (c *WorkerPool) Run(number int) {
	// no need to do http call by ourselves since omise has it's own API in golang, it abstract a lot of things to us and make our code more readable
	client, err := omise.NewClient(os.Getenv("OMISE_PUBLIC_KEY"), os.Getenv("OMISE_SECRET_KEY"))
	defer c.Wg.Done()
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"public_key": os.Getenv("OMISE_PUBLIC_KEY"),
			"secret_key": os.Getenv("OMISE_SECRET_KEY"),
		}).Error("Connection to Omise API failed")
		return
	}
	for {
		// this is called receive operator, it's a native mechanism to know if the channel is empty AND closed so we can quit gracefully
		// https://golang.org/ref/spec#Receive_operator
		row, ok := <-c.ChannelWork
		if !ok {
			break
		}
		for {
			if retry := c.charge(row, client); retry == false {
				break
			}
			log.Warn(fmt.Sprintf("Worker #%d Reached API requests limit, waiting 5s", number))
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *WorkerPool) charge(rowPtr *[]string, client *omise.Client) bool {
	// Creates a token from a test card.
	// here we don't actually check the errors on Atoi, if the values are not numbers, we will get 0
	// it will then fail later on, so it's not necessary to check
	row := *rowPtr
	log.Println(row)
	amount, _ := strconv.ParseInt(row[c.Positions.Amount], 10, 0)
	year, _ := strconv.Atoi(row[c.Positions.Year])
	month, _ := strconv.Atoi(row[c.Positions.Month])
	token, createToken := &omise.Token{}, &operations.CreateToken{
		Name:            row[c.Positions.Name],
		Number:          row[c.Positions.Card],
		ExpirationMonth: time.Month(month),
		ExpirationYear:  year,
	}
	if e := client.Do(token, createToken); e != nil {
		retry := true
		if checkLimitAPIError(e) == false {
			retry = false
			c.Stats.AddFailedTx(amount)
		}
		return retry
	}

	// Creates a charge from the token
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   amount, // ฿ 1,000.00
		Currency: os.Getenv("CURRENCY"),
		Card:     token.ID,
	}
	if e := client.Do(charge, createCharge); e != nil {
		retry := true
		if checkLimitAPIError(e) == false {
			retry = false
			c.Stats.AddFailedTx(amount)
		}
		return retry
	}
	c.Stats.AddSuccessTx(amount)
	return false
}

func checkLimitAPIError(err error) bool {
	log.Error(err)
	// we can't look for *omise.Error, because the limit is enforced via NGINX and therefore recognize as a ErrTransport
	if errTransport, ok := err.(*omise.ErrTransport); ok {
		return strings.Contains(errTransport.Error(), "429 Too Many Requests")
	}
	return false
}
