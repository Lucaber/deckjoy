# DeckJoy
Connect your Steam Deck as USB GameController/Keyboard/Mouse to your PC. No additional software or drivers required.

## Installation
### Requirements
- Password of the `deck` linux user 
  - Set a password can be set with the `passwd` command in a terminal
- Steam Deck in USB DRD Mode
  - Shutdown your Steam Deck
  - Hold `Volume UP` and press the `Power` button
  - Navigate to `Setup Utilities` using the D-Pad and confirm with A
  - Navigate to `Advanced` -> `USB Configuration` -> `USB Dual Role Device` using the D-Pad
  - Change the setting to `DRD`
  - Navigate to `Exit` -> `Exit Saving Changes` -> `Yes`

### Installation
- Switch to Desktop Mode
- Download and unpack the latest deckjoy.zip
- Add deckjoy as a Non-Steam Game 
- Switch back to Gaming Mode
- Start deckjoy

### Steam Input

If Steam Input is not configured automatically I recommend to configure:

- Trackpads
  - Right Trackpad Behaviour: As Mouse
    - Click: Left Mouse Click
- Action Sets
  - Default
    - Add Always-On Command
      - Command: System -> Touchscreen Native Support
