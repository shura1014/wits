package wits

import (
	"fmt"
	"github.com/shura1014/logger"
)

var (
	version = "1.0.0"
	banner  = `        .__  __          
__  _  _|__|/  |_  ______
\ \/ \/ /  \   __\/  ___/
 \     /|  ||  |  \___ \ 
  \/\_/ |__||__| /____  >
                      \/ %s
Hello, Welcome To Use Wits.
`
)

func SetBanner(b string) {
	banner = b
}

func printBanner() {
	if GetBool(AppBannerEnable, true) {
		fmt.Printf(logger.Cyan(banner)+"\n", version)
	}
}
