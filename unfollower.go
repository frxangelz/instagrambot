/* unfollower
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jimlawless/whereami"
)

func unfollower_main() {
	fmt.Println("Hello World!")
}

type _UNFOLLOWER struct {
	prev, next *_UNFOLLOWER
	id         int64
	date       int64
}

type cUnfollower struct {
	enable                  bool
	chance                  int
	min_node                int
	max_node                int
	search_interval         int64
	follow_back_expired     int64 // 0 = disable
	post_activity_expired   int64 // 0 = disable
	max_unfollow_in_session int
	max_unfollow_in_day     int
	unfollow_interval       int64
	debug                   bool
	exception_usernames     []string

	// counter
	last_search               int64
	last_unfollow             int64
	count                     int
	unfollow_count_in_session int
	unfollow_count_in_day     int

	last_max_id string
	// classes
	Followings *cInvData

	following_max_id_fname string
}

func (U *cUnfollower) Initialize() {

	U.following_max_id_fname = "./following_last_max_id.txt"
	section, err := conf.Section("unfollower")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		halt()
	}

	U.enable = section.ValueOf("enable") == "1"
	if !U.enable {
		mylog(whereami.WhereAmI(), "Unfollower is OFF")
		return
	}

	U.min_node = strToInt(section.ValueOf("min_node"))
	U.max_node = strToInt(section.ValueOf("max_node"))

	U.chance = strToInt(section.ValueOf("chance"))
	if (U.chance < 1) || (U.chance > 100) {
		U.chance = 75
		mylog(whereami.WhereAmI(), "invalid chance value, set to default 75 %")
	}

	U.search_interval = strToInt64(section.ValueOf("search_interval"))
	if U.search_interval < 30 {
		mylog(whereami.WhereAmI(), "search_interval too fast, set to default 30 minutes")
		U.search_interval = 30
	}
	U.search_interval = U.search_interval * 60

	U.follow_back_expired = strToInt64(section.ValueOf("follow_back_expired"))
	U.follow_back_expired = U.follow_back_expired * 3600 * 24
	U.post_activity_expired = strToInt64(section.ValueOf("post_activity_expired"))
	U.post_activity_expired = U.post_activity_expired * 3600 * 24

	U.max_unfollow_in_session = strToInt(section.ValueOf("max_unfollow_in_session"))
	U.max_unfollow_in_day = strToInt(section.ValueOf("max_unfollow_in_day"))
	U.unfollow_interval = strToInt64(section.ValueOf("unfollow_interval"))
	U.unfollow_interval = U.unfollow_interval * 60

	U.debug = section.ValueOf("debug") == "1"

	U.count = 0
	U.unfollow_count_in_session = 0
	U.unfollow_count_in_day = 0

	ts := strings.Split(section.ValueOf("exception_usernames"), ",")
	for _, s := range ts {

		if s != "" {
			U.exception_usernames = append(U.exception_usernames, s)
		}
	}

	if section.ValueOf("run_on_start") != "1" {

		U.last_search = time.Now().Unix()
	}

	// classes
	U.Followings = new(cInvData)
	U.Followings.Initialize()
	U.Followings.fname = "./followings.txt"
	U.Followings.load()
	U.Followings.ExpiredSecs = U.follow_back_expired
	U.Followings.autosave = true

	U._load_last_max_id()
}

func (U *cUnfollower) _load_last_max_id() {

	if !IsFileExists(U.following_max_id_fname) {
		U.last_max_id = ""
		return
	}

	f, err := os.Open(U.following_max_id_fname)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	U.last_max_id = string(b)
}

func (U *cUnfollower) _save_last_max_id() {

	f, err := os.Create(U.following_max_id_fname)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return
	}

	f.WriteString(U.last_max_id)
	f.Close()
}

func (U *cUnfollower) _is_exception_names(username string) bool {

	for _, s := range U.exception_usernames {

		if s == username {
			return true
		}
	}

	return false
}

func (U *cUnfollower) new_session() {

	U.unfollow_count_in_session = 0
}

func (U *cUnfollower) new_day() {

	U.unfollow_count_in_day = 0
}

func (U *cUnfollower) _unfollow_not_follow_back(CurTime int64) bool {

	if U.follow_back_expired == 0 {
		return false
	}

	if U.unfollow_count_in_session >= U.max_unfollow_in_session {
		return false
	}

	if U.unfollow_count_in_day >= U.max_unfollow_in_day {
		return false
	}

	if CurTime-U.last_unfollow < U.unfollow_interval {
		return false
	}

	res := U.Followings.GetOneExpired(CurTime)
	if res == nil {
		return false
	}

	U.last_unfollow = CurTime
	// get friendship status
	fr, err := Instagram.UserFriendShip(res.id)
	if err != nil {

		mylog(whereami.WhereAmI(), err)
		U.Followings.del(res.id)
		// hanya utk menandai, jangan lakukan process lain
		return true
	}

	if (fr.Following) && (!fr.FollowedBy) {

		U.unfollow_count_in_session++
		U.unfollow_count_in_day++

		_, err = Instagram.UnFollow(res.id)
		if err != nil {
			mylog(whereami.WhereAmI(), err)
		} else {
			s := intToStr(U.unfollow_count_in_session) + ":" + intToStr(U.max_unfollow_in_session) +
				" - " + intToStr(U.unfollow_count_in_day) + ":" + intToStr(U.max_unfollow_in_day)
			mylog(whereami.WhereAmI(), "Unfollowed (Not Followback):", res.id, s)
		}
	}

	// delete node
	U.Followings.del(res.id)
	return true
}

func (U *cUnfollower) _search(CurTime int64) bool {

	if U.count > U.min_node {
		// there're still node waiting for process, just dont search
		return false
	}

	if CurTime-U.last_search < U.search_interval {
		return false
	}

	if !getChance(U.chance) {
		return false
	}

	U.last_search = CurTime

	// searching unfollower
	user, err := Instagram.SelfUserFollowing(U.last_max_id)
	if err != nil {

		mylog(whereami.WhereAmI(), err)
		return true
	}

	if len(user.Users) < 1 {
		mylog(whereami.WhereAmI(), "No Following Last")
		U.last_max_id = ""
		return true
	}

	U.last_max_id = user.NextMaxID
	i := 0
	for _, u := range user.Users {

		if U._is_exception_names(u.Username) {
			continue
		}

		// check for auto follow
		b := false
		if AutoFollow != nil {

			b = AutoFollow.Followings.find(u.ID) != nil
		}

		if !b {
			if U.Followings.add(u.ID, CurTime, true) {
				i++
			}
		}

	}

	if U.Followings.modified {
		U.Followings.save()
	}

	U._save_last_max_id()
	if i > 0 {
		mylog(whereami.WhereAmI(), "Added Following :", i)
	}

	return true
}

func (U *cUnfollower) execute(CurTime int64) bool {

	if U._unfollow_not_follow_back(CurTime) {
		return true
	}

	return U._search(CurTime)
}
