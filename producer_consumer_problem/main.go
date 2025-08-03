package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberofBurgerRequired = 3

var burgersReady, burgersFailed, totalBurgers int

type Producer struct {
	data chan BurgerOrder //producer will get burger order and will make in just minutes. this will get the orders for the burger
	quit chan chan error  // channel stores the  reason for finishing pizza making. It tells us when we are done making pizzas
}

type BurgerOrder struct {
	burgerNumber int
	description  string
	success      bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func MakeBurger(burgerNumber int) *BurgerOrder {
	burgerNumber++

	// if the current number of burger is less than or equal to required number of  burger required then we make burger otherwise not
	if burgerNumber <= NumberofBurgerRequired {
		//making burger, it might get success sometime and might get failed sometime.

		// this will generate a number in the half open interval [0,n).We need at least 1 minute so we are adding 1 if it is 0.
		timeToMakeBurger := rand.Intn(10)
		if timeToMakeBurger == 0 {
			timeToMakeBurger += 1
		}

		fmt.Printf("Successfully Received order #%d!\n", burgerNumber)

		// now we will generated random values within certain interval and will divide the interval into success and failure reason
		rand := rand.Intn(14) + 1
		msg := ""
		success := false

		// if the random value is less than 4 then burger making is failed and will success if more than or equal to the 4
		if rand < 4 {
			burgersFailed++
		} else {
			burgersReady++
		}

		totalBurgers++

		fmt.Printf("Order #%d is being prepared. It will take %d minutes \n", burgerNumber, timeToMakeBurger)

		//delay for certain time to depict burger is being ready. here we have delay of minutes but will delay for corresponding seconds.
		// means if we need timeToMakeBurger 5 min then we will just wait 5 seconds
		time.Sleep(time.Duration(timeToMakeBurger) * time.Second)

		//there will be multiple reasons for failure
		if rand <= 2 {
			msg = fmt.Sprintf("**We  ran out of ingredients for burger #%d!", burgerNumber)
		} else if rand < 4 {
			msg = fmt.Sprintf("**Chef making burger #%d resigned", burgerNumber)
		} else {
			success = true
			msg = fmt.Sprintf("**Order for burger #%d is ready!", burgerNumber)

		}

		return &BurgerOrder{
			burgerNumber: burgerNumber,
			description:  msg,
			success:      success,
		}
	}

	return &BurgerOrder{
		burgerNumber: burgerNumber,
	}
}

func BurgerKitchen(burgerChef *Producer) {
	//todos:
	// keeping the track of which burger  is being made
	// run forever or until we receive quit  notification
	//trying to make burgers

	var burgerNumber = 0
	for {
		// try to make burgers
		// get the information about pizza like ready failed or ? based on the certain conditions
		currentBurger := MakeBurger(burgerNumber)

		if currentBurger != nil {
			burgerNumber = currentBurger.burgerNumber
			select {
			// we tried to make burger and send something to data channel
			case burgerChef.data <- *currentBurger:

			case quitChan := <-burgerChef.quit:
				//close both channels
				close(burgerChef.data)
				close(quitChan)
				return
			}
		}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	color.Cyan("Welcome to the Burger House !")
	color.Cyan("------------------------------------")

	//creating a producer

	burgerChef := Producer{
		data: make(chan BurgerOrder),
		quit: make(chan chan error),
	}

	//run the producer in background
	go BurgerKitchen(&burgerChef)

	//creating a consumer
	for i := range burgerChef.data {
		if i.burgerNumber <= NumberofBurgerRequired {
			if i.success {
				color.Green(i.description)
				color.Green("Order #%d is sent for delivery!", i.burgerNumber)
			} else {
				color.Red(i.description)
				color.Red("Failed to make burger")
			}
		} else {
			color.Cyan("Burger making is completed!")
			err := burgerChef.Close()
			if err != nil {
				color.Red("*** Failed to Close channel!")
			}
		}
	}

	color.Cyan("-------------")
	color.Cyan("Done.")

	color.Cyan("Among %d burgers, failed to make %d, and succeed %d burgers.", totalBurgers, burgersFailed, burgersReady)

}
