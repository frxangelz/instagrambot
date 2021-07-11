/* session
Instagram Bot - autofollow, like, comment dan unfollow not followback (c) Free Angel - frxangelz@gmail.com
please subscribe to my channel :
https://www.youtube.com/channel/UC15iFd0nlfG_tEBrt6Qz1NQ
*/
package main

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/ahmdrz/goinsta/store"
	"github.com/jimlawless/whereami"
)

// Logins and saves the session
func createAndSaveSession() error {
	err := Instagram.Login()
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return err
	}

	key := createKey()
	bytes, err := store.Export(Instagram, key)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return err
	}
	err = ioutil.WriteFile("session", bytes, 0644)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return err
	}

	return nil
}

// reloadSession will attempt to recover a previous session
func reloadSession() error {
	if _, err := os.Stat("session"); os.IsNotExist(err) {
		mylog(whereami.WhereAmI(), "No Session Found")
		return errors.New("No session found")
	}

	session, err := ioutil.ReadFile("session")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return err
	}

	mylog(whereami.WhereAmI(), "A session file exists")

	key, err := ioutil.ReadFile("key")
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return err
	}

	Instagram, err = store.Import(session, key)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return errors.New("Couldn't recover the session")
	}

	mylog(whereami.WhereAmI(), "Successfully logged in")
	return nil

}

// createKey creates a key and saves it to file
func createKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return nil
	}
	err = ioutil.WriteFile("key", key, 0644)
	if err != nil {
		mylog(whereami.WhereAmI(), err)
		return nil
	}

	//	mylog(whereami.WhereAmI(), "Created and saved the key", key)
	return key
}
