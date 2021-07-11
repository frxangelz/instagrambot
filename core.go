/* core
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ
*/

package main

import (
	"fmt"
	"time"

	"github.com/ahmdrz/goinsta"
	"github.com/jimlawless/whereami"
)

func core_main() {
	fmt.Println("Hello World!")
}

// jadikan global var agar bisa diakses dari semua class, initator di module core
var (
	Instagram  *goinsta.Instagram
	Unfollower *cUnfollower
	Feeds      *cFeeds
	AutoFollow *cAutoFollow
	Follower   *cFollower
)

type cCore struct {
	interval    time.Duration // in seconds
	CurTime     time.Time
	UnixTime    int64 // cur time in unix
	SessionTime int64 // in seconds, unix time

	CurrentDay int
	debug      bool
}

func (C *cCore) Initialize() {

	C.debug = false
	C.interval = 30 * time.Second
	C.CurTime = time.Now()
	C.UnixTime = C.CurTime.Unix()
	C.SessionTime = 3600 // default 1 hour
	C.CurrentDay = C.CurTime.Day()

	// unfollower
	Unfollower = new(cUnfollower)
	Unfollower.Initialize()

	Follower = new(cFollower)
	Follower.Initialize()

	// feeds
	Feeds = new(cFeeds)
	Feeds.Initialize()

	// autofollow
	AutoFollow = new(cAutoFollow)
	AutoFollow.Initialize()

	/* instagram =================================================== */
	section, err := conf.Section("main")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		halt()
	}

	user := section.ValueOf("username")
	password := section.ValueOf("password")
	Instagram = goinsta.New(user, password)

	err = reloadSession()
	if err != nil {
		err = createAndSaveSession()
		if err != nil {
			halt()
		}
	}
}

func (C *cCore) synch() {

	for {
		time.Sleep(time.Duration(config.synch_interval) * time.Second)
		Instagram.SyncFeatures()
	}
}

func (C *cCore) execute() {

	Session := 1
	_cur_time := C.UnixTime
	mylog(whereami.WhereAmI(), "Session :", Session)

	_is_new_session := false
	_last_session_time := _cur_time
	_is_new_day := false

	if config.synch_interval > 0 {
		go C.synch()
	}

	for {
		if C.debug {
			fmt.Println("sleep ...")
		}

		time.Sleep(C.interval)

		C.CurTime = time.Now()
		_cur_time = C.CurTime.Unix()
		_is_new_session = _cur_time-_last_session_time > C.SessionTime
		d := C.CurTime.Day()
		_is_new_day = d != C.CurrentDay

		C.UnixTime = C.CurTime.Unix()
		C.CurrentDay = d

		// check for session
		if _is_new_session {

			// new session
			Session++
			mylog(whereami.WhereAmI(), "Session :", Session)
			_last_session_time = _cur_time

			// do whatever to do in new section
			if Unfollower != nil {
				Unfollower.new_session()
			}

			if AutoFollow != nil {
				AutoFollow.new_session()
			}

			if Follower != nil {
				Follower.new_session()
			}
		}

		if _is_new_day {

			mylog(whereami.WhereAmI(), "day :", C.CurrentDay)
			// do whatever to do in new day
			if Unfollower != nil {
				Unfollower.new_day()
			}

			if AutoFollow != nil {
				AutoFollow.new_day()
			}

			if Follower != nil {
				Follower.new_day()
			}
		}

		if Feeds != nil {

			Feeds.execute(_cur_time)
		}

		// unfollower
		if Unfollower != nil {

			if Unfollower.execute(_cur_time) {

				if C.debug {
					fmt.Println("Unfollower return success ...")
				}
				//continue
			}
		}

		if Follower != nil {

			if Follower.execute(_cur_time) {
				if C.debug {
					fmt.Println("Follower return success ...")
				}
				//continue
			}
		}

		if AutoFollow != nil {
			if AutoFollow.execute(_cur_time) {
				if C.debug {
					fmt.Println("Autofollow return success ...")
				}
				//continue
			}
		}

	}
}
