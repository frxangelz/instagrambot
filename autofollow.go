/* autofollow
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ

- need feeds.go
*/

package main

import (
	"fmt"
	"strings"

	/*	"io/ioutil"
		"os" */
	"time"

	"github.com/jimlawless/whereami"
)

func autofollow_main() {
	fmt.Println("Hello World!")
}

type cAutoFollow struct {
	enable bool
	chance int

	like_interval       int64
	max_like_in_session int
	max_like_in_day     int

	max_follow_in_session   int
	max_follow_in_day       int
	follow_interval         int64
	max_unfollow_in_session int
	max_unfollow_in_day     int
	follow_expired          int64
	unfollow_interval       int64
	debug                   bool
	exception_usernames     []string
	comments                []string
	comment_chance          int
	// counter
	//	last_search               int64
	last_follow             int64
	last_like               int64
	count                   int
	follow_count_in_session int
	follow_count_in_day     int
	like_count_in_session   int
	like_count_in_day       int

	last_unfollow             int64
	unfollow_count_in_session int
	unfollow_count_in_day     int

	last_max_id string

	// classes
	Followings *cInvData
}

func (A *cAutoFollow) Initialize() {

	section, err := conf.Section("auto")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		halt()
	}

	A.enable = section.ValueOf("enable") == "1"
	if !A.enable {
		mylog(whereami.WhereAmI(), "Auto is OFF")
		return
	}

	A.chance = strToInt(section.ValueOf("chance"))
	if (A.chance < 1) || (A.chance > 100) {
		A.chance = 75
		mylog(whereami.WhereAmI(), "invalid chance value, set to default 75 %")
	}

	A.like_interval = strToInt64(section.ValueOf("like_interval"))
	if A.like_interval < 1 {
		mylog(whereami.WhereAmI(), "like_interval too fast, set to default 5 minutes")
		A.like_interval = 5
	}
	A.like_interval = A.like_interval * 60
	A.max_like_in_session = strToInt(section.ValueOf("max_like_in_session"))
	A.max_like_in_day = strToInt(section.ValueOf("max_like_in_day"))

	A.max_follow_in_session = strToInt(section.ValueOf("max_follow_in_session"))
	A.max_follow_in_day = strToInt(section.ValueOf("max_follow_in_day"))
	A.follow_interval = strToInt64(section.ValueOf("follow_interval"))
	if A.follow_interval < 1 {
		A.follow_interval = 5
		mylog(whereami.WhereAmI(), "follow_interval too fast, set to default 5 minutes")
	}
	A.follow_interval = A.follow_interval * 60

	A.max_unfollow_in_session = strToInt(section.ValueOf("max_unfollow_in_session"))
	A.max_unfollow_in_day = strToInt(section.ValueOf("max_unfollow_in_day"))
	A.unfollow_interval = strToInt64(section.ValueOf("unfollow_interval"))
	if A.unfollow_interval < 3 {
		A.unfollow_interval = 5
		mylog(whereami.WhereAmI(), "unfollow_interval too fast, set to default 5 minutes")
	}
	A.unfollow_interval = A.unfollow_interval * 60

	A.follow_expired = strToint64(section.ValueOf("follow_expired"))
	if A.follow_expired < 1 {
		A.follow_expired = 30
		mylog(whereami.WhereAmI(), "invalid follow_expired value, set to default 30 days")
	}
	A.follow_expired = A.follow_expired * 24 * 3600

	ts := strings.Split(section.ValueOf("exception_usernames"), ",")
	for _, s := range ts {
		if s != "" {
			A.exception_usernames = append(A.exception_usernames, s)
		}
	}
	ts = strings.Split(section.ValueOf("comment"), "|")
	for _, s := range ts {
		if s != "" {
			A.comments = append(A.comments, s)
		}
	}
	A.comment_chance = strToInt(section.ValueOf("comment_chance"))

	A.debug = section.ValueOf("debug") == "1"

	if section.ValueOf("run_on_start") != "1" {
		A.last_follow = time.Now().Unix()
		A.last_like = A.last_follow
		A.last_unfollow = A.last_follow
	}

	// classes
	A.Followings = new(cInvData)
	A.Followings.Initialize()
	A.Followings.fname = "./autofollow.txt"
	A.Followings.load()
	A.Followings.ExpiredSecs = A.follow_expired
	A.Followings.autosave = true

	if A.debug {
		fmt.Println(whereami.WhereAmI(), "max_like ", A.max_like_in_session, ":", A.max_like_in_day)
		fmt.Println(whereami.WhereAmI(), "max_follow ", A.max_follow_in_session, ":", A.max_follow_in_day)
		fmt.Println(whereami.WhereAmI(), "max_unfollow ", A.max_unfollow_in_session, ":", A.max_unfollow_in_day)
	}
}

func (A *cAutoFollow) _is_exception_username(username string) bool {

	for i := 0; i < len(A.exception_usernames); i++ {
		if strings.Contains(username, A.exception_usernames[i]) {
			return true
		}
	}

	return false
}

func (A *cAutoFollow) _get_comment() string {

	if A.comment_chance == 0 {
		return ""
	}

	if len(A.comments) == 0 {
		return ""
	}

	if !getChance(A.comment_chance) {
		return ""
	}

	i := random(0, len(A.comments)-1)
	return A.comments[i]
}

func (A *cAutoFollow) _follow(CurTime int64) bool {

	if CurTime-A.last_follow < A.follow_interval {
		return false
	}

	if A.follow_count_in_session >= A.max_follow_in_session {
		return false
	}

	if A.follow_count_in_day >= A.max_follow_in_day {
		return false
	}

	A.last_follow = CurTime

	if !getChance(A.chance) {
		return false
	}

	feed := Feeds.GetForFollow()
	if feed.id == 0 {
		return false
	}

	A.follow_count_in_day++
	A.follow_count_in_session++

	_, err := Instagram.Follow(feed.user_id)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return true
	}

	s := intToStr(A.follow_count_in_session) + ":" + intToStr(A.max_follow_in_session) + " - " +
		intToStr(A.follow_count_in_day) + ":" + intToStr(A.max_follow_in_day)

	if A.Followings.add(feed.user_id, CurTime, true) {
		if A.Followings.autosave {
			A.Followings.save()
		}
	}

	mylog(whereami.WhereAmI(), "followed UserId :", feed.user_id, s)
	return true
}

func (A *cAutoFollow) _like(CurTime int64) bool {

	if CurTime-A.last_like < A.like_interval {
		return false
	}

	if A.like_count_in_session >= A.max_like_in_session {
		return false
	}

	if A.like_count_in_day >= A.max_like_in_day {
		return false
	}

	A.last_like = CurTime

	if !getChance(A.chance) {
		if A.debug {
			mylog(whereami.WhereAmI(), A.chance, "Not My Lucky Day, ouch !")
		}
		return false
	}

	feed := Feeds.GetForLike()
	if feed.id == 0 {
		if A.debug {
			mylog(whereami.WhereAmI(), "GetForLike return 0")
		}
		return false
	}

	A.like_count_in_day++
	A.like_count_in_session++

	_, err := Instagram.Like(uint64ToStr(uint64(feed.id)))
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return true
	}

	s := intToStr(A.like_count_in_session) + ":" + intToStr(A.max_like_in_session) + " - " +
		intToStr(A.like_count_in_day) + ":" + intToStr(A.max_like_in_day)

	mylog(whereami.WhereAmI(), "liked Media Id :", feed.id, s)

	s = A._get_comment()
	if s != "" {
		Instagram.Comment(uint64ToStr(uint64(feed.id)), s)
		mylog(whereami.WhereAmI(), "commented :", feed.id, s)
	}

	return true
}

func (A *cAutoFollow) _unfollow(CurTime int64) bool {

	if CurTime-A.last_unfollow < A.unfollow_interval {
		return false
	}

	A.last_unfollow = CurTime

	if !getChance(A.chance) {
		return false
	}

	// get expired
	user := A.Followings.GetOneExpired(CurTime)
	if user == nil {
		if A.debug {
			mylog(whereami.WhereAmI(), "No User Expired")
		}
		return false
	}

	A.unfollow_count_in_day++
	A.unfollow_count_in_session++

	_, err := Instagram.UnFollow(user.id)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
	} else {

		s := intToStr(A.unfollow_count_in_session) + ":" + intToStr(A.max_unfollow_in_session) + " - " +
			intToStr(A.unfollow_count_in_day) + ":" + intToStr(A.max_unfollow_in_day)

		mylog(whereami.WhereAmI(), "Unfollowed (Not Followback) :", user.id, s)
	}

	// delete
	A.Followings.del(user.id)
	return true
}

func (A *cAutoFollow) execute(CurTime int64) bool {

	if A._like(CurTime) {
		return true
	}

	if A._follow(CurTime) {
		return true
	}

	if A._unfollow(CurTime) {
		return true
	}

	return false
}

func (A *cAutoFollow) new_session() {

	A.follow_count_in_session = 0
	A.like_count_in_session = 0
	A.unfollow_count_in_session = 0
}

func (A *cAutoFollow) new_day() {
	A.follow_count_in_day = 0
	A.like_count_in_day = 0
	A.unfollow_count_in_day = 0
}
