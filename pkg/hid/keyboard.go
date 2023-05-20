package hid

var KeyboardReportDesc = []byte{
	0x05, 0x01, /* USAGE_PAGE (Generic Desktop)	          */
	0x09, 0x06, /* USAGE (Keyboard)                       */
	0xa1, 0x01, /* COLLECTION (Application)               */
	0x05, 0x07, /*   USAGE_PAGE (Keyboard)                */
	0x19, 0xe0, /*   USAGE_MINIMUM (Keyboard LeftControl) */
	0x29, 0xe7, /*   USAGE_MAXIMUM (Keyboard Right GUI)   */
	0x15, 0x00, /*   LOGICAL_MINIMUM (0)                  */
	0x25, 0x01, /*   LOGICAL_MAXIMUM (1)                  */
	0x75, 0x01, /*   REPORT_SIZE (1)                      */
	0x95, 0x08, /*   REPORT_COUNT (8)                     */
	0x81, 0x02, /*   INPUT (Data,Var,Abs)                 */
	//0x95, 0x01, /*   REPORT_COUNT (1)                     */
	//0x75, 0x08, /*   REPORT_SIZE (8)                      */
	//0x81, 0x03, /*   INPUT (Cnst,Var,Abs)                 */
	//0x95, 0x05, /*   REPORT_COUNT (5)                     */
	//0x75, 0x01, /*   REPORT_SIZE (1)                      */
	//0x05, 0x08, /*   USAGE_PAGE (LEDs)                    */
	//0x19, 0x01, /*   USAGE_MINIMUM (Num Lock)             */
	//0x29, 0x05, /*   USAGE_MAXIMUM (Kana)                 */
	//0x91, 0x02, /*   OUTPUT (Data,Var,Abs)                */
	//0x95, 0x01, /*   REPORT_COUNT (1)                     */
	//0x75, 0x03, /*   REPORT_SIZE (3)                      */
	//0x91, 0x03, /*   OUTPUT (Cnst,Var,Abs)                */
	0x95, 0x06, /*   REPORT_COUNT (6)                     */
	0x75, 0x08, /*   REPORT_SIZE (8)                      */
	0x15, 0x00, /*   LOGICAL_MINIMUM (0)                  */
	0x25, 0x65, /*   LOGICAL_MAXIMUM (101)                */
	0x05, 0x07, /*   USAGE_PAGE (Keyboard)                */
	0x19, 0x00, /*   USAGE_MINIMUM (Reserved)             */
	0x29, 0x65, /*   USAGE_MAXIMUM (Keyboard Application) */
	0x81, 0x00, /*   INPUT (Data,Ary,Abs)                 */
	0xc0, /* END_COLLECTION                         */
}

type Keyboard struct {
	*Device
	modKeysPressed map[KeyboardModKey]any
	keysPressed    map[KeyboardKey]any
}

func (k *Keyboard) PressMod(key KeyboardModKey) error {
	k.modKeysPressed[key] = struct{}{}
	return k.SendState()
}
func (k *Keyboard) ReleaseMod(key KeyboardModKey) error {
	delete(k.modKeysPressed, key)
	return k.SendState()
}
func (k *Keyboard) Press(key KeyboardKey) error {
	k.keysPressed[key] = struct{}{}
	return k.SendState()
}
func (k *Keyboard) Release(key KeyboardKey) error {
	delete(k.keysPressed, key)
	return k.SendState()
}

func (k *Keyboard) SendState() error {
	state := []byte{0x00}

	for modKey := range k.modKeysPressed {
		state[0] |= byte(modKey)
	}

	if len(k.keysPressed) <= 6 {
		for key := range k.keysPressed {
			state = append(state, byte(key))
		}
		for len(state) < 7 {
			state = append(state, 0x00)
		}
	} else {
		for len(state) < 7 {
			state = append(state, 0x01)
		}
	}

	return k.Write(state)
}

func NewKeyboard(path string) *Keyboard {
	k := &Keyboard{
		Device: &Device{
			path: path,
		},
		keysPressed:    map[KeyboardKey]any{},
		modKeysPressed: map[KeyboardModKey]any{},
	}

	return k
}
