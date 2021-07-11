/* content
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ

- follower, load to memory, delete node after action eg like,mention etc
*/
package main

import (
	"fmt"
	"time"

	"github.com/jimlawless/whereami"
)

type _FOLLOWER_DATA struct {
	prev, next         *_FOLLOWER_DATA
	id                 int64
	username           string
	mentioned          bool
	followback_checked bool
}

type cFollower struct {
	enable   bool
	chance   int
	interval int64
	mention  int

	max_follow_back_in_day     int
	max_follow_back_in_session int
	follow_back_interval       int64

	//run_on_start=1
	debug bool

	// counter
	last_search               int64
	last_follow_back          int64
	session_follow_back_count int
	day_follow_back_count     int

	head     *_FOLLOWER_DATA
	count    int
	max_node int

	max_id string
}

func (F *cFollower) _find(id int64) *_FOLLOWER_DATA {

	p := F.head.next
	for {
		if p == F.head {
			return nil
		}

		if p.id == id {
			return p
		}

		p = p.next
	}
}

func (F *cFollower) _add(id int64, username string) bool {

	if F.max_node == 0 {
		return false
	}

	p := F._find(id)
	if p != nil {
		// already exists
		return false
	}

	if F.count >= F.max_node {
		// delete first node
		first := F.head.next
		F.head.next = first.next
		first.next.prev = F.head

		first.username = ""
		first.next = nil
		first.prev = nil
		first = nil
		F.count--
	}

	// add last
	p = new(_FOLLOWER_DATA)
	p.next = F.head
	p.prev = F.head.prev
	F.head.prev.next = p
	F.head.prev = p

	p.id = id
	p.mentioned = false
	p.followback_checked = false
	p.username = username
	F.count++
	return true
}

func (F *cFollower) del(id int64) {

	p := F._find(id)
	if p == nil {
		return
	}

	p.prev.next = p.next
	p.next.prev = p.prev
	p.username = ""
	p.next = nil
	p.prev = nil
	p = nil
	F.count--
}

func (F *cFollower) walk() {

	p := F.head.next
	for {
		if p == F.head {
			break
		}

		fmt.Println(p.id, p.username)
		p = p.next
	}
}
func (F *cFollower) Initialize() {

	F.head = new(_FOLLOWER_DATA)
	F.head.next = F.head
	F.head.prev = F.head
	F.max_node = 5000
	F.count = 0

	section, err := conf.Section("follower")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		halt()
	}

	F.enable = section.ValueOf("enable") == "1"
	if !F.enable {
		mylog(whereami.WhereAmI(), "Follower is OFF")
		return
	}

	F.max_node = strToInt(section.ValueOf("max_node"))

	F.chance = strToInt(section.ValueOf("chance"))
	if (F.chance < 1) || (F.chance > 100) {
		mylog(whereami.WhereAmI(), "invalid chance value, set to default 75")
		F.chance = 75
	}

	F.interval = strToInt64(section.ValueOf("interval"))
	if F.interval < 10 {
		F.interval = 15
		mylog(whereami.WhereAmI(), "invalid interval value, set to default 15")
	}
	F.interval = F.interval * 60

	F.mention = strToInt(section.ValueOf("mention"))

	F.max_follow_back_in_day = strToInt(section.ValueOf("max_follow_back_in_day"))
	F.max_follow_back_in_session = strToInt(section.ValueOf("max_follow_back_in_session"))
	F.follow_back_interval = strToInt64(section.ValueOf("follow_back_interval"))
	if F.follow_back_interval < 3 {
		mylog(whereami.WhereAmI(), "follow_back_interval too fast, set to default 3 minutes")
		F.follow_back_interval = 3
	}
	F.follow_back_interval = F.follow_back_interval * 60

	F.debug = section.ValueOf("debug") == "1"
	if section.ValueOf("run_on_start") != "1" {
		F.last_search = time.Now().Unix()
		F.last_follow_back = F.last_search
	}

	F.max_id = ""
}

//func (F)
func (F *cFollower) _get_name_to_mention() string {

	if F.head.next == F.head {
		return ""
	}

	p := F.head.next
	for {

		if p == F.head {
			return ""
		}

		if !p.mentioned {
			p.mentioned = true
			return p.username
		}
		p = p.next
	}

	return ""
}

func (F *cFollower) GetMentions() (res string) {

	if F.mention < 1 {
		return ""
	}

	if !getChance(F.chance) {
		return ""
	}

	res = ""
	i := 0
	for {
		if i >= F.mention {
			return res
		}

		s := F._get_name_to_mention()
		if s == "" {
			return res
		}

		i++
		if i == 1 {
			res = "@" + s
		} else {
			res = res + " @" + s
		}
	}
}

// return one data used to delte
func (F *cFollower) _get_trash() *_FOLLOWER_DATA {

	FollowbackEnable := (F.max_follow_back_in_day > 0) && (F.max_follow_back_in_session > 0)

	p := F.head.next
	for {
		if p == F.head {
			return nil
		}

		if !FollowbackEnable {

			if p.mentioned {

				return p
			}
		} else {

			if p.mentioned && p.followback_checked {
				return p
			}
		}

		p = p.next
	}

	return nil
}

func (F *cFollower) ClearTrash() {

	for {
		p := F._get_trash()
		if p == nil {
			return
		}

		p.prev.next = p.next
		p.next.prev = p.prev
		p.username = ""
		p.next = nil
		p.prev = nil
		p = nil
		F.count--
	}
}

func (F *cFollower) _get_follow_back_not_checked() *_FOLLOWER_DATA {

	p := F.head.next
	for {
		if p == F.head {
			return nil
		}

		if !p.followback_checked {
			return p
		}

		p = p.next
	}

	return nil
}

func (F *cFollower) _follow_back(CurTime int64) bool {

	if F.session_follow_back_count >= F.max_follow_back_in_session {
		return false
	}

	if F.day_follow_back_count >= F.max_follow_back_in_day {
		return false
	}

	if CurTime-F.last_follow_back < F.follow_back_interval {
		return false
	}

	F.last_follow_back = CurTime
	if !getChance(F.chance) {
		return false
	}

	// check for followback
	p := F._get_follow_back_not_checked()
	if p == nil {
		return false
	}

	res, err := Instagram.UserFriendShip(p.id)
	p.followback_checked = true
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return true
	}

	if !res.FollowedBy {

		F.day_follow_back_count++
		F.session_follow_back_count++

		_, err = Instagram.Follow(p.id)
		if err != nil {
			mylog(whereami.WhereAmI(), err)
		} else {
			s := intToStr(F.session_follow_back_count) + ":" + intToStr(F.max_follow_back_in_session) + "," +
				intToStr(F.day_follow_back_count) + ":" + intToStr(F.max_follow_back_in_day)
			mylog(whereami.WhereAmI(), "Followback :", p.username, s)
			s = ""
		}
	}

	return true
}

func (F *cFollower) _search(CurTime int64) bool {

	if CurTime-F.last_search < F.interval {
		return false
	}

	F.last_search = CurTime
	if !getChance(F.chance) {
		return false
	}

	res, err := Instagram.SelfUserFollowers(F.max_id)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return true
	}

	if len(res.Users) < 1 {
		F.max_id = ""
		return true
	}

	i := 0
	for _, user := range res.Users {

		if F._add(user.ID, user.Username) {
			i++
		}
	}

	if i > 0 {
		mylog(whereami.WhereAmI(), "Added To Node :", i)
	}
	return true
}

func (F *cFollower) execute(CurTime int64) bool {

	if !F.enable {
		return false
	}

	if F._search(CurTime) {
		return true
	}

	if F._follow_back(CurTime) {
		return true
	}

	F.ClearTrash()
	return false
}

func (F *cFollower) new_session() {

	F.session_follow_back_count = 0
}

func (F *cFollower) new_day() {
	F.day_follow_back_count = 0
}
