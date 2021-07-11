/* feeds
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ
*/
package main

import (
	"fmt"
	"strings"

	"github.com/jimlawless/whereami"
)

func feeds_main() {
	fmt.Println("Hello World!")
}

type _FEED_DATA struct {
	prev, next       *_FEED_DATA
	id               int64 // media id
	takenAt          int64 // datetime
	caption          string
	user_id          int64
	user_is_favorite bool
	user_is_private  bool
}

type _FEED_TAG_DATA struct {
	tag       string
	use_count int
}

type cFeeds struct {
	enable   bool
	chance   int
	interval int64
	tags     []_FEED_TAG_DATA

	expired  int64 // int seconds (unix time), config in days
	min_like int   // if likes < min_like will skip
	debug    bool

	//node
	head     *_FEED_DATA
	count    int
	max_node int

	// counter
	last_search int64
}

func (F *cFeeds) Initialize() {

	F.head = new(_FEED_DATA)
	F.head.next = F.head
	F.head.prev = F.head
	F.count = 0
	F.max_node = 100

	section, err := conf.Section("feeds")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		halt()
	}

	F.enable = section.ValueOf("enable") == "1"
	if !F.enable {
		mylog(whereami.WhereAmI(), "feeds is OFF")
		return
	}

	F.chance = strToInt(section.ValueOf("chance"))
	if (F.chance < 1) || (F.chance > 100) {
		mylog(whereami.WhereAmI(), "invalid chance value, set to default 75 %")
		return
	}

	F.max_node = strToInt(section.ValueOf("max_node"))
	F.interval = strToInt64(section.ValueOf("interval"))
	if F.interval < 10 {
		mylog(whereami.WhereAmI(), "interval too fast, set to default 60 (1 hours)")
		F.interval = 60
	}
	F.interval = F.interval * 60

	ts := strings.Split(section.ValueOf("tags"), ",")
	var ftd _FEED_TAG_DATA
	ftd.use_count = 0

	for _, s := range ts {
		if s != "" {
			ftd.tag = s
			F.tags = append(F.tags, ftd)
		}
	}
	ts = nil

	F.expired = strToInt64(section.ValueOf("expired"))
	F.expired = F.expired * 24 * 3600

	F.min_like = strToInt(section.ValueOf("min_like"))
	F.debug = section.ValueOf("debug") == "1"
}

func (F *cFeeds) _find_feed(feed_id int64) *_FEED_DATA {

	p := F.head.next
	for {

		if p == F.head {
			return nil
		}

		if p.id == feed_id {
			return p
		}

		p = p.next
	}

	return nil
}

func (F *cFeeds) _find_feed_user(user_id int64) *_FEED_DATA {

	p := F.head.next
	for {

		if p == F.head {
			return nil
		}

		if p.user_id == user_id {
			return p
		}

		p = p.next
	}

	return nil
}

func (F *cFeeds) _add(id int64, title string, user_id int64, is_favorite, is_private bool, takenAt int64) {

	// only add different user
	if F._find_feed_user(user_id) != nil {
		return
	}

	if F.count >= F.max_node {
		return
	}

	p := new(_FEED_DATA)
	p.next = F.head
	p.prev = F.head.prev
	F.head.prev.next = p
	F.head.prev = p

	p.id = id
	p.caption = title
	p.user_id = user_id
	p.user_is_favorite = is_favorite
	p.user_is_private = is_private
	p.takenAt = takenAt
	F.count++
}

func (F *cFeeds) _del(id int64) {

	p := F._find_feed(id)
	if p == nil {
		return
	}

	p.prev.next = p.next
	p.next.prev = p.prev
	p.caption = ""
	p.next = nil
	p.prev = nil
	F.count--
}

func (F *cFeeds) _get_tag_to_search() string {

	if len(F.tags) < 1 {
		return ""
	}

	if len(F.tags) == 1 {
		return F.tags[0].tag
	}

	use_count := F.tags[0].use_count
	idx := 0

	for i, s := range F.tags {

		if s.use_count < use_count {
			use_count = s.use_count
			idx = i
		}
	}

	F.tags[idx].use_count++
	return F.tags[idx].tag
}

func (F *cFeeds) _get_popular_feeds(CurTime int64) {

	res, err := Instagram.GetPopularFeed()
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return
	}

	if F.debug {
		mylog(whereami.WhereAmI(), "Get Popular Feeds return :", res.NumResults)
	}

	if len(res.Items) < 1 {
		return
	}

	for _, item := range res.Items {

		// check for liked
		if item.HasLiked {
			continue
		}

		if item.LikeCount < F.min_like {
			continue
		}

		if CurTime-item.TakenAt > F.expired {
			continue
		}

		//Instagram.Like()
		F._add(item.Caption.MediaID, item.Caption.Text, item.User.Pk, item.User.IsFavorite, item.User.IsPrivate, item.TakenAt)
	}
}

func (F *cFeeds) _search_by_tag(CurTime int64, tag string) {

	res, err := Instagram.TagFeed(tag)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return
	}

	if F.debug {
		mylog(whereami.WhereAmI(), "Search Tag Feeds", tag, ", found :", res.NumResults)
	}

	if len(res.Items) < 1 {
		return
	}

	for _, item := range res.Items {

		// check for liked
		if item.HasLiked {
			continue
		}

		if item.LikeCount < F.min_like {
			continue
		}

		if CurTime-item.TakenAt > F.expired {
			continue
		}
		//Instagram.Like()
		F._add(item.Caption.MediaID, item.Caption.Text, item.User.ID, item.User.IsFavorite, item.User.IsPrivate, item.TakenAt)
	}
}

func (F *cFeeds) _get_one_expired(CurTime int64) *_FEED_DATA {

	p := F.head.next
	for {
		if p == F.head {
			return nil
		}

		if CurTime-p.takenAt > F.expired {
			return p
		}

		p = p.next
	}

	return nil
}
func (F *cFeeds) ClearExpired(CurTime int64) {

	p := F.head.next
	for {

		p = F._get_one_expired(CurTime)
		if p == nil {
			return
		}

		F._del(p.id)
	}

}

func (F *cFeeds) GetForLike() (res _FEED_DATA) {

	p := F.head.next
	res.id = 0
	if p == F.head {
		return res
	}

	res = *p
	p.prev.next = p.next
	p.next.prev = p.prev
	p.caption = ""
	p.prev = nil
	p.next = nil
	p = nil
	F.count--
	return res
}

func (F *cFeeds) GetForFollow() (res _FEED_DATA) {
	p := F.head.next
	res.id = 0

	for {
		if p == F.head {
			return res
		}

		if p.user_is_favorite {
			p = p.next
			continue
		}

		if p.user_is_private {
			p = p.next
			continue
		}

		break
	}

	res = *p
	p.prev.next = p.next
	p.next.prev = p.prev
	p.caption = ""
	p.prev = nil
	p.next = nil
	p = nil
	F.count--
	return res
}

func (F *cFeeds) _search(CurTime int64) {

	if CurTime-F.last_search < F.interval {
		return
	}

	if F.count > F.max_node {
		return
	}

	F.last_search = CurTime

	if !getChance(F.chance) {
		return
	}

	tag := F._get_tag_to_search()
	if tag == "" {
		F._get_popular_feeds(CurTime)
	} else {

		F._search_by_tag(CurTime, tag)
	}
}

func (F *cFeeds) execute(CurTime int64) {

	if !F.enable {
		return
	}

	F.ClearExpired(CurTime)
	F._search(CurTime)
}
