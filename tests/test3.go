// Copyright 2018 Alexander S.Kresin <alex@kresin.ru>, http://www.kresin.ru
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
// A sample from a "The Go programming language", a book of Alan A. A. Donovan & Brian W. Kernighan,
// adapted by Alexander S.Kresin for External GUI framework.
package main

import (
	egui "github.com/alkresin/external"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math/cmplx"
	"os"
)

const (
	CLR_LBLUE  = 16759929
	CLR_LBLUE0 = 12164479
	CLR_LBLUE2 = 16770002
	CLR_LBLUE3 = 16772062
	CLR_LBLUE4 = 16775920

	CLR_LGRAY1 = 15658734
	CLR_LGRAY2 = 14540253
)

var width, height = 680, 680

func main() {

	var sInit string

	{
		b, err := ioutil.ReadFile("test.ini")
		if err != nil {
			sInit = ""
		} else {
			sInit = string(b)
		}
	}

	if !egui.Init(sInit) {
		return
	}

	egui.CreateStyle(&(egui.Style{Name: "st1", Orient: 1, Colors: []int32{CLR_LBLUE, CLR_LBLUE3}}))
	egui.CreateStyle(&(egui.Style{Name: "st2", Colors: []int32{CLR_LBLUE}, BorderW: 3}))
	egui.CreateStyle(&(egui.Style{Name: "st3", Colors: []int32{CLR_LBLUE},
		BorderW: 2, BorderClr: CLR_LBLUE0}))
	pFont := egui.CreateFont(&(egui.Font{Name: "f1", Family: "Georgia", Height: -14}))

	pWindow := &(egui.Widget{X: 100, Y: 100, W: 716, H: 764, Title: "Test3", BColor: 1, Font: pFont})
	egui.InitMainWindow(pWindow)

	pPanel := pWindow.AddWidget(&(egui.Widget{Type: "paneltop", H: 32,
		AProps: map[string]string{"HStyle": "st1"}}))

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 0, Y: 0, W: 60, H: 32, Title: "Exit",
		AProps: map[string]string{"HStyles": egui.ToString("st1", "st2", "st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", nil, "hwg_EndWindow()")

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 60, Y: 0, W: 60, H: 32, Title: "M-brot",
		AProps: map[string]string{"HStyles": egui.ToString("st1", "st2", "st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", fu1, "fu1", "1")

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 120, Y: 0, W: 60, H: 32, Title: "Acos",
		AProps: map[string]string{"HStyles": egui.ToString("st1", "st2", "st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", fu1, "fu1", "2")

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 180, Y: 0, W: 60, H: 32, Title: "Sqrt",
		AProps: map[string]string{"HStyles": egui.ToString("st1", "st2", "st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", fu1, "fu1", "3")

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 240, Y: 0, W: 60, H: 32, Title: "Newton",
		AProps: map[string]string{"HStyles": egui.ToString("st1", "st2", "st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", fu1, "fu1", "4")

	pWindow.AddWidget(&(egui.Widget{Type: "bitmap", Name: "img", X: 10, Y: 36, W: 680,
		H: 680, BColor: CLR_LGRAY2}))
	//	Anchor: egui.A_TOPABS+egui.A_LEFTABS+egui.A_BOTTOMABS+egui.A_RIGHTABS }))
	//egui.PLastWidget.SetCallBackProc("onsize", nil, "{|o,x,y|o:Move(,,x-o:nLeft,y-72)}")

	pWindow.Activate()

	egui.Exit()

}

func fu1(p []string) string {

	egui.PLastWindow.Move(-1, -1, 716, 764)
	pImg := egui.Widg("main.img")
	pImg.SetImage("")

	switch p[1] {
	case "1":
		draw(mandelbrot)
	case "2":
		draw(acos)
	case "3":
		draw(sqrt)
	case "4":
		draw(newton)
	}

	pImg.SetImage("a1.jpg")
	return ""
}

func draw(fu func(z complex128) color.Color) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/float64(height)*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/float64(width)*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, fu(z))
		}
	}
	fo, _ := os.Create("a1.jpg")
	jpeg.Encode(fo, img, nil) // NOTE: ignoring errors
	fo.Close()
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}

func acos(z complex128) color.Color {
	v := cmplx.Acos(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}

func sqrt(z complex128) color.Color {
	v := cmplx.Sqrt(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{128, blue, red}
}

// f(x) = x^4 - 1
//
// z' = z - f(z)/f'(z)
//    = z - (z^4 - 1) / (4 * z^3)
//    = z - (z - 1/z^3) / 4
func newton(z complex128) color.Color {
	const iterations = 37
	const contrast = 7
	for i := uint8(0); i < iterations; i++ {
		z -= (z - 1/(z*z*z)) / 4
		if cmplx.Abs(z*z*z*z-1) < 1e-6 {
			return color.Gray{255 - contrast*i}
		}
	}
	return color.Black
}
