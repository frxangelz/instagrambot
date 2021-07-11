/*invdata
 content
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

	"github.com/jimlawless/whereami"
)

type _INV_DATA struct {
	prev, next *_INV_DATA
	id         int64
	date       int64
}

type cInvData struct {
	max_node int
	count    int
	debug    bool
	modified bool
	autosave bool
	fname    string
	head     *_INV_DATA

	ExpiredSecs int64 // expired dalam detik
}

func (I *cInvData) clear() {

	p := I.head.next
	for {
		if p == I.head {
			break
		}

		p1 := p.next
		p.next = nil
		p.prev = nil
		p = nil
		p = p1
	}

	I.head.next = I.head
	I.head.prev = I.head
	I.count = 0
	I.modified = true
}

func (I *cInvData) Initialize() {

	I.head = new(_INV_DATA)
	I.head.next = I.head
	I.head.prev = I.head
	I.count = 0
	I.max_node = 15000
	I.modified = false
	I.ExpiredSecs = 3600 * 24 * 30 // 30 hari default
	I.debug = true
	I.autosave = true

	I.fname = "./inv.txt"
}

func (I *cInvData) find(id int64) *_INV_DATA {

	p := I.head.next
	for {
		if p == I.head {
			return nil
		}

		if p.id == id {
			return p
		}

		p = p.next
	}

	return nil
}

func (I *cInvData) add(id, date int64, check bool) bool {

	var p *_INV_DATA
	if check {

		p = I.find(id)
		if p != nil {
			return false
		}
	}

	if I.count >= I.max_node {
		if I.debug {
			mylog(whereami.WhereAmI(), "Max Node :", I.count, ":", I.max_node)
		}
		return false
	}

	p = new(_INV_DATA)
	p.next = I.head
	p.prev = I.head.prev
	I.head.prev.next = p
	I.head.prev = p
	I.count++

	p.id = id
	p.date = date
	/*
		if I.debug {
			mylog(whereami.WhereAmI(), "New Id :", p.id, I.count, ":", I.max_node)
		}
	*/
	I.modified = true
	return true
}

func (I *cInvData) GetOneExpired(CurTime int64) *_INV_DATA {

	p := I.head.next
	for {
		if p == I.head {
			return nil
		}

		if CurTime-p.date > I.ExpiredSecs {
			return p
		}

		p = p.next
	}

	return nil
}

func (I *cInvData) ClearExpired(CurTime int64) {

	for {

		p := I.GetOneExpired(CurTime)
		if p == nil {
			return
		}

		// delete node
		p.next.prev = p.prev
		p.prev.next = p.next
		p.next = nil
		p.prev = nil
		I.count--
		if I.debug {
			mylog(whereami.WhereAmI(), "deleted :", p.id)
		}
		p = nil
		I.modified = true
	}
}

func (I *cInvData) save() bool {

	if !I.modified {
		return false
	}

	p := I.head.next
	if p == I.head {

		return false
	}

	f, _ := os.Create(I.fname)

	for {

		if p == I.head {
			break
		}

		s := int64ToStr(p.id) + "," + int64ToStr(p.date) + "\n"
		_, err := f.WriteString(s)
		if err != nil {
			mylog(whereami.WhereAmI(), err)
			break
		}

		p = p.next
	}

	I.modified = false
	return true
}

func (I *cInvData) walk() {
	p := I.head.next

	for {
		if p == I.head {
			break
		}

		fmt.Println(p.id, p.date)
		p = p.next
	}
}

func (I *cInvData) load() bool {

	if !IsFileExists(I.fname) {
		return false
	}

	f, err := os.Open(I.fname)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return false
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)

	ts := strings.Split(string(b), "\n")
	b = nil
	var id int64 = 0
	var date int64 = 0

	for _, s := range ts {

		sx := strings.Split(s, ",")
		if len(sx) < 2 {
			//mylog(whereami.WhereAmI(), "invalid string value", s)
			break
		}
		id = strToInt64(sx[0])
		if id == 0 {
			//mylog(whereami.WhereAmI(), "invalid string value", s)
			break
		}
		date = strToint64(sx[1])
		I.add(id, date, false)
	}

	mylog(whereami.WhereAmI(), "Loaded :", I.fname, I.count, ":", I.max_node)
	I.modified = false
	return true
}

func (I *cInvData) IsExists(id int64) bool {
	return I.find(id) != nil
}

func (I *cInvData) del(id int64) {

	p := I.find(id)
	if p == nil {
		return
	}

	p.next.prev = p.prev
	p.prev.next = p.next
	p.next = nil
	p.prev = nil
	p = nil
	I.count--

	if I.autosave {
		I.save()
	}
}
