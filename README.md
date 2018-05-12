# Tello Drone Control

Very early version with a lot of hacks.

**Needless to say I take no responsibility for any damage you caused to your drone. You should know exactly what you are doing before attempting to use any of this code. Be extremely cautious with gamepad control as it might work entirely differently with your gamepad.**

What works:
* Basic controls using a gamepad (up, down, rotate left, rotate right, forward, backward, left, right)
* Video stream
* Flip controls work most of the time

## Controller mappings

These are meant for Dualshock 4 controllers. Yours might be completely different.

### Left Analog
- Up - Ascend
- Down - Descend
- Left - Rotate Counter-clockwise
- Right - Rotate Clockwise

### Right Analog
- Up - Forward
- Down - Backward
- Left - Sideways left
- Right - Sideways right

### D-pad
- Up - Flip forward
- Back - Flip backward
- Left - Flip left
- Right - Flip right

### Buttons
- X - Take off / Land (toggle)

## Preqrequisites

* libavcodec - Used for decoding H.264 packets

## Installation

```
go get github.com/socketbind/drone-control
$GOPATH/bin/drone-control
```
