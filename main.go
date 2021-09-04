/* main
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ
*/
package main

import (
	"fmt"
	"os"

	"github.com/jimlawless/whereami"
)

func halt() {

	CleanUp()
	if Instagram != nil {
		Instagram.Logout()
	}
	os.Exit(0)
}

func main() {
	fmt.Println("Instagram Bot - autofollow, like, comment dan unfollow not followback (c) free Angel")
	fmt.Println("please subscribe to my youtube channel (Newbie Computer) :\n https://www.youtube.com/channel/UCqRqvw9n7Lrh79x3dRDOkDg")
	fmt.Println("----------------------------------------\n")

	Initialize()
	err := load_config(mypath + "config.txt")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return
	}

	err = CreateDir("tmp")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return
	}

	core := new(cCore)
	core.Initialize()
	core.execute()
}
