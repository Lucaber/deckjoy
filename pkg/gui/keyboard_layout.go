package gui

import (
	"github.com/lucaber/deckjoy/pkg/hid"
	"github.com/veandco/go-sdl2/sdl"
)

type Key struct {
	WidthUnits float64
	ModKey     hid.KeyboardModKey
	Key        hid.KeyboardKey
	Hidden     bool
	Text       string

	renderX       int32
	renderY       int32
	renderW       int32
	renderH       int32
	renderSurface *sdl.Surface
	pressed       bool
}

type KeyRow struct {
	// Left to Right
	Keys []Key
}

type Keyboard struct {
	// Top to Bottom
	Rows        []KeyRow
	PreRendered bool
}

var KeyboardAnsiTKL = Keyboard{
	Rows: []KeyRow{
		{
			Keys: []Key{
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyESC,
					Text:       "esc",
				},
				{
					WidthUnits: 1,
					Hidden:     true,
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF1,
					Text:       "F1",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF2,
					Text:       "F2",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF3,
					Text:       "F3",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF4,
					Text:       "F4",
				},
				{
					WidthUnits: 0.5,
					Hidden:     true,
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF5,
					Text:       "F5",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF6,
					Text:       "F6",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF7,
					Text:       "F7",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF8,
					Text:       "F8",
				},
				{
					WidthUnits: 0.5,
					Hidden:     true,
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF9,
					Text:       "F9",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF10,
					Text:       "F10",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF11,
					Text:       "F11",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF12,
					Text:       "F12",
				},
				{
					WidthUnits: 0.5,
					Hidden:     true,
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeySYSRQ,
					Text:       "sysrq",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeySCROLLLOCK,
					Text:       "scroll",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyPAUSE,
					Text:       "pause",
				},
			},
		},
		{
			Keys: []Key{
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyGRAVE,
					Text:       "` ~",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey1,
					Text:       "1 !",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey2,
					Text:       "2 @",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey3,
					Text:       "3 #",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey4,
					Text:       "4 $",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey5,
					Text:       "5 %",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey6,
					Text:       "6 ^",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey7,
					Text:       "7 &",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey8,
					Text:       "8 *",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey9,
					Text:       "9 (",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKey0,
					Text:       "0 )",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyMINUS,
					Text:       "- _",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyEQUAL,
					Text:       "= +",
				},
				{
					WidthUnits: 2,
					Key:        hid.KeyboardKeyBACKSPACE,
					Text:       "backspace",
				},
				{
					WidthUnits: 0.5,
					Hidden:     true,
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyINSERT,
					Text:       "ins",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyHOME,
					Text:       "home",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyPAGEUP,
					Text:       "pgup",
				},
			},
		},
		{
			Keys: []Key{
				{
					WidthUnits: 1.5,
					Key:        hid.KeyboardKeyTAB,
					Text:       "tab",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyQ,
					Text:       "Q",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyW,
					Text:       "W",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyE,
					Text:       "E",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyR,
					Text:       "R",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyT,
					Text:       "T",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyY,
					Text:       "Y",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyU,
					Text:       "U",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyI,
					Text:       "I",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyO,
					Text:       "O",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyP,
					Text:       "P",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyLEFTBRACE,
					Text:       "[ {",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyRIGHTBRACE,
					Text:       "] }",
				},
				{
					WidthUnits: 1.5,
					Key:        hid.KeyboardKeyBACKSLASH,
					Text:       "\\ |",
				},
				{
					WidthUnits: 0.5,
					Hidden:     true,
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyDELETE,
					Text:       "del",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyEND,
					Text:       "end",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyPAGEDOWN,
					Text:       "pgdn",
				},
			},
		},
		{
			Keys: []Key{
				{
					WidthUnits: 1.75,
					Key:        hid.KeyboardKeyCAPSLOCK,
					Text:       "caps",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyA,
					Text:       "A",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyS,
					Text:       "S",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyD,
					Text:       "D",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyF,
					Text:       "F",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyG,
					Text:       "G",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyH,
					Text:       "H",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyJ,
					Text:       "J",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyK,
					Text:       "K",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyL,
					Text:       "L",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeySEMICOLON,
					Text:       "; :",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyAPOSTROPHE,
					Text:       "' \"",
				},
				{
					WidthUnits: 2.25,
					Key:        hid.KeyboardKeyENTER,
					Text:       "return",
				},
			},
		},
		{
			Keys: []Key{
				{
					WidthUnits: 2.25,
					ModKey:     hid.KeyboardModKeyLShift,
					Text:       "shift",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyZ,
					Text:       "Z",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyX,
					Text:       "X",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyC,
					Text:       "C",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyV,
					Text:       "V",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyB,
					Text:       "B",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyN,
					Text:       "N",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyM,
					Text:       "M",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyCOMMA,
					Text:       ", <",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyDOT,
					Text:       ". >",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeySLASH,
					Text:       "/ ?",
				},
				{
					WidthUnits: 2.75,
					ModKey:     hid.KeyboardModKeyRShift,
					Text:       "shift",
				},
				{
					WidthUnits: 1.5,
					Hidden:     true,
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyUP,
					Text:       "^",
				},
			},
		},
		{
			Keys: []Key{
				{
					WidthUnits: 1.25,
					ModKey:     hid.KeyboardModKeyLCtrl,
					Text:       "ctrl",
				},
				{
					WidthUnits: 1.25,
					ModKey:     hid.KeyboardModKeyLMeta,
					Text:       "meta",
				},
				{
					WidthUnits: 1.25,
					ModKey:     hid.KeyboardModKeyLAlt,
					Text:       "alt",
				},
				{
					WidthUnits: 6.25,
					Key:        hid.KeyboardKeySPACE,
					Text:       "space",
				},
				{
					WidthUnits: 1.25,
					ModKey:     hid.KeyboardModKeyRAlt,
					Text:       "alt",
				},
				{
					WidthUnits: 1.25,
					ModKey:     hid.KeyboardModKeyRMeta,
					Text:       "meta",
				},
				{
					WidthUnits: 1.25,
					Key:        hid.KeyboardKeyCOMPOSE,
					Text:       "menu",
				},
				{
					WidthUnits: 1.25,
					ModKey:     hid.KeyboardModKeyRCtrl,
					Text:       "ctrl",
				},
				{
					WidthUnits: 0.5,
					Hidden:     true,
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyLEFT,
					Text:       "<",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyDOWN,
					Text:       "v",
				},
				{
					WidthUnits: 1,
					Key:        hid.KeyboardKeyRIGHT,
					Text:       ">",
				},
			},
		},
	},
}
