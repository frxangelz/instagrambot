/* cinstagram
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ
*/

package main

import (
	"fmt"
	"time"

	"github.com/ahmdrz/goinsta"
)

func cinstagram_main() {
	fmt.Println("Hello World!")
}

type cInstagram struct {
	enable        bool
	interval      int64
	photo_quality int // default 87

	last_post int64
}

func (I *cInstagram) Initialize() {

	I.last_post = time.Now().Unix()
}

func test() {

	insta := goinsta.New("myjendhil", "cilacap2008")

	if err := insta.Login(); err != nil {
		panic(err)
	}

	defer insta.Logout()

	fmt.Println("Logged In")

	//insta.Explore()
	id := insta.NewUploadID()

	//res, err := insta.UploadVideo(mypath+"tmp/aS2PUzxGkihPJydE.mp4", mypath+"tmp/DYL3JJvVoAAGjar.jpg", "funny people", id)
	res, err := insta.UploadPhoto(mypath+"tmp/broke.jpg", "Broke", id, 86, 0)
	if err != nil {
		fmt.Println(err)
		// instagram error
	} else {
		fmt.Println("video posted :", id, "funny people", res.Status)
	}
	//res, err := insta.Like("1734295949923487387") //_3416684")
	//
	/*	res, err := insta.Comment("1734295949923487387", "nice !!!")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(res))
		}
	*/
	/*	user, err := insta.SelfUserFollowing("")
			if err != nil {
				fmt.Println(err)
			}

		fmt.Println("NextMaxID :", user.NextMaxID, "count :", len(user.Users))
		for _, u := range user.Users {
			fmt.Println(u.ID, u.Username)
		}
	*/

	//res, _ := insta.SearchTags("#viral")
	/*	res, _ := //insta.TagFeed("funny")
				insta.SelfUserFollowers("")
				//insta.Timeline("")
			for _, item := range res.Users {
				//fmt.Println(item.Name, item.ID, item.MediaCount)
				//fmt.Println(item.Caption.MediaID, time.Unix(item.TakenAt, 0), item.HasLiked, item.LikeCount)
				fmt.Println(item.ID, item.Username, item.IsFavorite)
			}

		n := time.Now().Add(-180 * 24 * time.Hour)
		res, _ := insta.UserFeed(1990837127, "", int64ToStr(n.Unix()))
		for _, item := range res.Items {
			fmt.Println(item.Caption.Text)
		}
	*/
}
