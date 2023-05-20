package steamworks

import (
	"fmt"
	"github.com/ebitengine/purego"
	"os"
	"time"
)

type SteamInput uintptr
type InputHandle uintptr
type ActionSetHandle uintptr

var SteamAPIInit func() bool
var GetSteamInput func() SteamInput
var SteamInputInit func(SteamInput, bool) bool
var SetInputActionManifestFilePath func(SteamInput, string) bool
var GetActionSetHandle func(SteamInput, string) ActionSetHandle

var RunCallbacks func()
var RunFrame func(SteamInput, bool)

var GetConnectedControllers func(SteamInput, *uint64) int
var GetControllerForGamepadIndex func(SteamInput, int) InputHandle
var GetCurrentActionSet func(SteamInput, InputHandle) ActionSetHandle
var ActivateActionSet func(SteamInput, InputHandle, ActionSetHandle)

var SteamInputInstance SteamInput

var DefaultActionSet ActionSetHandle
var TouchActionSet ActionSetHandle

func Init() error {
	path := "./libsteam_api.so"
	lib, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", path, err)
	}

	// var RestartAppIfNecessary func(uint322 uint32)
	// purego.RegisterLibFunc(&RestartAppIfNecessary, lib, "SteamAPI_RestartAppIfNecessary")

	purego.RegisterLibFunc(&SteamAPIInit, lib, "SteamAPI_Init")
	purego.RegisterLibFunc(&GetSteamInput, lib, "SteamAPI_SteamInput_v006")
	purego.RegisterLibFunc(&SteamInputInit, lib, "SteamAPI_ISteamInput_Init")
	purego.RegisterLibFunc(&SetInputActionManifestFilePath, lib, "SteamAPI_ISteamInput_SetInputActionManifestFilePath")
	purego.RegisterLibFunc(&GetActionSetHandle, lib, "SteamAPI_ISteamInput_GetActionSetHandle")

	purego.RegisterLibFunc(&RunCallbacks, lib, "SteamAPI_RunCallbacks")
	purego.RegisterLibFunc(&RunFrame, lib, "SteamAPI_ISteamInput_RunFrame")

	purego.RegisterLibFunc(&GetConnectedControllers, lib, "SteamAPI_ISteamInput_GetConnectedControllers")
	purego.RegisterLibFunc(&GetControllerForGamepadIndex, lib, "SteamAPI_ISteamInput_GetControllerForGamepadIndex")
	purego.RegisterLibFunc(&GetCurrentActionSet, lib, "SteamAPI_ISteamInput_GetCurrentActionSet")
	purego.RegisterLibFunc(&ActivateActionSet, lib, "SteamAPI_ISteamInput_ActivateActionSet")

	if !SteamAPIInit() {
		return fmt.Errorf("SteamAPI_Init failed")
	}

	SteamInputInstance = GetSteamInput()
	if SteamInputInstance == 0 {
		return fmt.Errorf("failed to get SteamInput")
	}

	err = setup()
	if err != nil {
		SteamInputInstance = 0
		return err
	}

	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			RunCallbacks()
			RunFrame(SteamInputInstance, false)
		}
	}()

	return nil
}

func setup() error {
	if !SteamInputInit(SteamInputInstance, false) {
		return fmt.Errorf("SteamInput_Init failed")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	manifestPath := fmt.Sprintf("%s/game_actions_480.vdf", cwd)
	if !SetInputActionManifestFilePath(SteamInputInstance, manifestPath) {
		return fmt.Errorf("SetInputActionManifestFilePath(%s) failed", manifestPath)
	}

	DefaultActionSet = GetActionSetHandle(SteamInputInstance, "Default")
	if DefaultActionSet == 0 {
		return fmt.Errorf("failed to get ActionSetHandle Default")
	}
	TouchActionSet = GetActionSetHandle(SteamInputInstance, "Touch")
	if TouchActionSet == 0 {
		return fmt.Errorf("failed to get ActionSetHandle Touch")
	}

	return nil
}

func GetController() {
	c := make([]uint64, 20, 20)
	_ = GetConnectedControllers(SteamInputInstance, &(c[0]))
}

func ActivateActionSetForDefaultController(actionSet ActionSetHandle) error {
	if SteamInputInstance == 0 {
		return fmt.Errorf("steamworks not initialized")
	}

	input := GetControllerForGamepadIndex(SteamInputInstance, 0)
	if input == 0 {
		return fmt.Errorf("failed to get default controller")
	}

	current := GetCurrentActionSet(SteamInputInstance, input)
	if current == actionSet {
		return nil
	}
	ActivateActionSet(SteamInputInstance, input, actionSet)
	return nil
}
