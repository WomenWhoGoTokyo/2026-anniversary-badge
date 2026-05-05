package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"machine"
	"math"
	"sync/atomic"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/st7789"
)

// 240x240 RGB565
//go:embed images/badge.rgb565
var badgeImage []byte

// ロータリーエンコーダー 押下
var buzzerEnabled atomic.Bool

// 上下キースイッチ 押下
var (
	upperPressed atomic.Bool
	lowerPressed atomic.Bool
)

// ジョイスティック 傾き
var (
	joystickEngaged atomic.Bool
	joystickSector  atomic.Int32
)

// ディスプレイの回転角 (度) 0..359
var displayAngle atomic.Int32

// 1 クリックあたりの回転角 (度)
const rotationStepDeg int32 = 15

// 回転描画の 1 行分バッファ
var rotatedRow [240 * 2]byte

// SK6812 RGBW
type rgbw struct {
	R, G, B, W uint8
}

func (c rgbw) encode() uint32 {
	return uint32(c.G)<<24 | uint32(c.R)<<16 | uint32(c.B)<<8 | uint32(c.W)
}

func mixRGBW(x, y rgbw) rgbw {
	return rgbw{
		R: uint8((uint16(x.R) + uint16(y.R)) / 2),
		G: uint8((uint16(x.G) + uint16(y.G)) / 2),
		B: uint8((uint16(x.B) + uint16(y.B)) / 2),
		W: uint8((uint16(x.W) + uint16(y.W)) / 2),
	}
}

func approachU8(cur, target, step uint8) uint8 {
	if cur < target {
		if target-cur < step {
			return target
		}
		return cur + step
	}
	if cur-target < step {
		return target
	}
	return cur - step
}

var (
	pinkColor = rgbw{R: 120, G: 40, B: 70, W: 30} // 桜色
	blueColor = rgbw{R: 25, G: 75, B: 120, W: 30} // 水色
)

var rainbowColors = []rgbw{
	{R: 130, G: 50, B: 60, W: 30},  // 0: 淡い赤
	{R: 130, G: 80, B: 30, W: 40},  // 1: 淡い橙
	{R: 120, G: 110, B: 30, W: 40}, // 2: 淡い黄
	{R: 50, G: 130, B: 60, W: 40},  // 3: 淡い緑
	{R: 30, G: 110, B: 120, W: 40}, // 4: 淡い水色
	{R: 40, G: 70, B: 130, W: 30},  // 5: 淡い青
	{R: 100, G: 50, B: 130, W: 30}, // 6: 淡い紫
}

type WS2812B struct {
	Pin machine.Pin
	ws  *piolib.WS2812B
}

func NewWS2812B(pin machine.Pin) *WS2812B {
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, pin)
	ws.EnableDMA(true)
	return &WS2812B{ws: ws}
}

func (ws *WS2812B) WriteRaw(rawGRB []uint32) error {
	return ws.ws.WriteRaw(rawGRB)
}

var (
	buzzerPWM = machine.PWM4
	buzzerCh  uint8
)

func initBuzzer() {
	buzzerPWM.Configure(
		machine.PWMConfig{
			Period: 1_000_000_000 / 440,
		},
	)
	ch, _ := buzzerPWM.Channel(machine.D8)
	buzzerCh = ch
	buzzerPWM.Set(buzzerCh, 0) // 無音
}

func playTone(freq uint16, duration time.Duration) {
	buzzerPWM.Configure(
		machine.PWMConfig{
			Period: uint64(1_000_000_000) / uint64(freq),
		},
	)
	buzzerPWM.Set(buzzerCh, buzzerPWM.Top()/2)
	end := time.Now().Add(duration)
	for time.Now().Before(end) {
		if !buzzerEnabled.Load() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	buzzerPWM.Set(buzzerCh, 0)
}

func runBlink() {
	led := machine.D7
	led.Configure(
		machine.PinConfig{
			Mode: machine.PinOutput,
		},
	)
	for {
		led.Low()
		time.Sleep(time.Millisecond * 1000)
		led.High()
		time.Sleep(time.Millisecond * 1000)
	}
}

func runBuzzer() {
	initBuzzer()

	// ドレミファソラシド
	notes := []uint16{262, 294, 330, 349, 392, 440, 494, 523}

	for {
		if !buzzerEnabled.Load() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		for _, n := range notes {
			if !buzzerEnabled.Load() {
				break
			}
			playTone(n, 300*time.Millisecond)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func runEncoderSW() {
	button1 := machine.D5
	button1.Configure(
		machine.PinConfig{
			Mode: machine.PinInputPullup,
		},
	)

	prev := true
	for {
		cur := button1.Get()
		if prev && !cur {
			next := !buzzerEnabled.Load()
			buzzerEnabled.Store(next)
			if next {
				println("buzzer ON")
			} else {
				println("buzzer OFF")
			}
		}
		prev = cur
		time.Sleep(20 * time.Millisecond)
	}
}

func runJoystick() {
	ax := machine.ADC{Pin: machine.GPIO29}
	ax.Configure(machine.ADCConfig{})
	ay := machine.ADC{Pin: machine.GPIO28}
	ay.Configure(machine.ADCConfig{})

	const printThreshold uint16 = 0x0400

	const engageDist int64 = 0x2000
	const engageDist2 = engageDist * engageDist

	cx := ax.Get()
	cy := ay.Get()
	fmt.Printf("joystick: %04X %04X (init)\n", cx, cy)

	lastX := cx
	lastY := cy
	for {
		x := ax.Get()
		y := ay.Get()

		ox := int64(x) - int64(cx)
		oy := int64(y) - int64(cy)
		if ox*ox+oy*oy > engageDist2 {
			angle := math.Atan2(float64(oy), float64(ox))
			n := len(rainbowColors)
			norm := (angle + math.Pi) / (2 * math.Pi)
			sector := int32(norm*float64(n)) % int32(n)
			if sector < 0 {
				sector += int32(n)
			}
			joystickSector.Store(sector)
			joystickEngaged.Store(true)
		} else {
			joystickEngaged.Store(false)
		}

		if absDiffU16(x, lastX) > printThreshold || absDiffU16(y, lastY) > printThreshold {
			fmt.Printf("joystick: %04X %04X\n", x, y)
			lastX = x
			lastY = y
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func absDiffU16(a, b uint16) uint16 {
	if a > b {
		return a - b
	}
	return b - a
}

func runKeyInput() {
	machine.GPIO0.Configure(
		machine.PinConfig{
			Mode: machine.PinInputPullup,
		},
	)
	machine.GPIO1.Configure(
		machine.PinConfig{
			Mode: machine.PinInputPullup,
		},
	)

	prevUp := false
	prevDown := false
	for {
		up := !machine.GPIO0.Get()
		down := !machine.GPIO1.Get()
		upperPressed.Store(up)
		lowerPressed.Store(down)

		if up && !prevUp {
			println("The top button was pressed")
		}
		if down && !prevDown {
			println("The bottom button was pressed")
		}
		prevUp = up
		prevDown = down

		time.Sleep(20 * time.Millisecond)
	}
}

func runRotary() {
	enc := encoders.NewQuadratureViaInterrupt(
		machine.GPIO2,
		machine.GPIO3,
	)
	enc.Configure(
		encoders.QuadratureConfig{
			Precision: 4,
		},
	)
	oldValue := 0
	for {
		if newValue := enc.Position(); newValue != oldValue {
			delta := int32(newValue - oldValue)
			a := displayAngle.Load() + delta*rotationStepDeg
			a = ((a % 360) + 360) % 360
			displayAngle.Store(a)
			println("rotary:", newValue, "angle:", a)
			oldValue = newValue
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func runLEDs() {
	ws := NewWS2812B(machine.GPIO4)

	const fadeStep uint8 = 3
	const tick = 20 * time.Millisecond

	var cur rgbw
	for {
		var target rgbw
		switch {
		case joystickEngaged.Load():
			idx := int(joystickSector.Load())
			if idx >= 0 && idx < len(rainbowColors) {
				target = rainbowColors[idx]
			}
		case upperPressed.Load() && lowerPressed.Load():
			target = mixRGBW(pinkColor, blueColor)
		case upperPressed.Load():
			target = pinkColor
		case lowerPressed.Load():
			target = blueColor
		}

		cur.R = approachU8(cur.R, target.R, fadeStep)
		cur.G = approachU8(cur.G, target.G, fadeStep)
		cur.B = approachU8(cur.B, target.B, fadeStep)
		cur.W = approachU8(cur.W, target.W, fadeStep)

		v := cur.encode()
		ws.WriteRaw([]uint32{v, v})

		time.Sleep(tick)
	}
}

var display st7789.Device

func setupDisplay() {
	machine.SPI1.Configure(
		machine.SPIConfig{
			Frequency: 16000000,
			Mode:      0,
		},
	)
	display = st7789.New(
		machine.SPI1,
		machine.GPIO9,
		machine.GPIO12,
		machine.GPIO13,
		machine.GPIO14,
	)
	display.Configure(
		st7789.Config{
			Height: 240,
			Width:  240,
		},
	)
	display.FillScreen(color.RGBA{0, 0, 0, 255})
}

func renderImage(angle int32) {
	if angle == 0 {
		display.DrawRGBBitmap8(0, 0, badgeImage, 240, 240)
		return
	}
	rad := float64(angle) * (math.Pi / 180)
	cosA := int32(math.Cos(rad) * 65536)
	sinA := int32(math.Sin(rad) * 65536)
	const cx, cy int32 = 120, 120
	const fp = 16

	for y := int32(0); y < 240; y++ {
		dy := y - cy
		sxFP := (cx << fp) + cosA*(-cx) + sinA*dy
		syFP := (cy << fp) - sinA*(-cx) + cosA*dy
		for x := int32(0); x < 240; x++ {
			sx := sxFP >> fp
			sy := syFP >> fp
			di := x * 2
			if uint32(sx) < 240 && uint32(sy) < 240 {
				si := (sy*240 + sx) * 2
				rotatedRow[di] = badgeImage[si]
				rotatedRow[di+1] = badgeImage[si+1]
			} else {
				rotatedRow[di] = 0
				rotatedRow[di+1] = 0
			}
			sxFP += cosA
			syFP -= sinA
		}
		display.DrawRGBBitmap8(0, int16(y), rotatedRow[:], 240, 1)
	}
}

func runDisplay() {
	var rendered int32 = 0
	for {
		cur := displayAngle.Load()
		if cur != rendered {
			renderImage(cur)
			rendered = cur
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func main() {
	machine.InitADC()

	setupDisplay()
	renderImage(0)

	go runBlink()
	go runBuzzer()
	go runEncoderSW()
	go runJoystick()
	go runKeyInput()
	go runRotary()
	go runLEDs()
	go runDisplay()

	// main を終わらせない
	select {}
}
