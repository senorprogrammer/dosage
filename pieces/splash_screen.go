package pieces

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
)

const splashMsg string = `
 ______   _____  _______ _______  ______ _______
 |     \ |     | |______ |_____| |  ____ |______
 |_____/ |_____| ______| |     | |_____| |______
                                                
`

const (
	splashInterval = 0
)

func DisplaySplashScreen() {
	fmt.Println(aurora.BrightGreen(splashMsg))
	time.Sleep(splashInterval * time.Second)
}