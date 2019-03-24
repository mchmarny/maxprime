package main

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/maxprime/util"
)

const (
	// maxNumber to check for primality - increase for a longer test
	defaultMaxNumber = 500000
)

// defaultPrimeHandler processes default prime calculator
func defaultPrimeHandler(c *gin.Context) {
	r := &resp{
		ID:      util.GetUUIDv4(),
		Prime:   *getPrimeResp(defaultMaxNumber),
		Ts:      time.Now().UTC().String(),
		Release: release,
	}
	c.JSON(http.StatusOK, r)
}

// primeWithArgHandler processes prime calculator requests with argument
func primeWithArgHandler(c *gin.Context) {

	maxVar := c.Param("max")
	log.Printf("max == %s", maxVar)
	if maxVar == "" {
		log.Println("Error on nil max parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Internal Error",
			"status":  http.StatusBadRequest,
		})
		return
	}

	maxNum, err := strconv.Atoi(maxVar)
	if err != nil {
		log.Printf("Error while parsing max parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Internal Error",
			"status":  http.StatusBadRequest,
		})
		return
	}

	r := &resp{
		ID:      util.GetUUIDv4(),
		Prime:   *getPrimeResp(maxNum),
		Ts:      time.Now().UTC().String(),
		Release: release,
	}
	c.JSON(http.StatusOK, r)

}

func getPrimeResp(maxNumber int) *prime {

	var x, y, n int
	nsqrt := math.Sqrt(float64(maxNumber))

	isPrime := make([]bool, maxNumber)

	for x = 1; float64(x) <= nsqrt; x++ {
		for y = 1; float64(y) <= nsqrt; y++ {
			n = 4*(x*x) + y*y
			if n <= maxNumber && (n%12 == 1 || n%12 == 5) {
				isPrime[n] = !isPrime[n]
			}
			n = 3*(x*x) + y*y
			if n <= maxNumber && n%12 == 7 {
				isPrime[n] = !isPrime[n]
			}
			n = 3*(x*x) - y*y
			if x > y && n <= maxNumber && n%12 == 11 {
				isPrime[n] = !isPrime[n]
			}
		}
	}

	for n = 5; float64(n) <= nsqrt; n++ {
		if isPrime[n] {
			for y = n * n; y < maxNumber; y += n * n {
				isPrime[y] = false
			}
		}
	}

	isPrime[2] = true
	isPrime[3] = true

	primes := make([]int, 0, 1270606)
	for x = 0; x < len(isPrime)-1; x++ {
		if isPrime[x] {
			primes = append(primes, x)
		}
	}

	return &prime{
		Max:   maxNumber,
		Value: primes[len(primes)-1],
	}
}

type resp struct {
	Prime   prime  `json:"prime"`
	ID      string `json:"id"`
	Ts      string `json:"ts"`
	Release string `json:"rel"`
}

type prime struct {
	Max   int `json:"max"`
	Value int `json:"val"`
}
