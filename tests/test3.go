// Copyright 2018 Alexander S.Kresin <alex@kresin.ru>, http://www.kresin.ru
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
// A sample from a "The Go programming language", a book of Alan A. A. Donovan & Brian W. Kernighan,
// adapted by Alexander S.Kresin for External GUI framework.
package main

import (
	"io/ioutil"
	"image"
	"image/color"
	"image/jpeg"
	"math/cmplx"
	"os"
	egui "github.com/alkresin/external"
)

const (
	CLR_LGRAY1  = 15658734
	CLR_LGRAY2  = 14540253
)

var width, height = 720, 720

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

	pWindow := &(egui.Widget{X: 100, Y: 100, W: 748, H: 800, Title: "Test3", BColor: 1})
	egui.InitMainWindow(pWindow)

	egui.Menu("")
	{
		egui.AddMenuItem("Exit", nil, "hwg_EndWindow()")
		egui.AddMenuItem("Mandelbrot", fu1, "fu1")
		egui.AddMenuItem("Acos", fu2, "fu2")
		egui.AddMenuItem("Sqrt", fu3, "fu3")
		egui.AddMenuItem("Newton", fu4, "fu4")
	}
	egui.EndMenu()

	//pWindow.AddWidget(&(egui.Widget{Type: "label", Name: "l1",
	//	X: 10, Y: 10, W: 180, H: 24, Title: ""}))

	pWindow.AddWidget(&(egui.Widget{Type: "bitmap", Name: "img", X: 10, Y: 10, W: 720,
		H: 720, BColor: CLR_LGRAY2 }))
	//	Anchor: egui.A_TOPABS+egui.A_LEFTABS+egui.A_BOTTOMABS+egui.A_RIGHTABS }))
	//egui.PLastWidget.SetCallBackProc("onsize", nil, "{|o,x,y|o:Move(,,x-o:nLeft,y-72)}")

	pWindow.Activate()

	egui.Exit()

}

func fu1([]string) string {

	draw(mandelbrot)
	
	pImg := egui.GetWidg("main.img")
	pImg.SetImage( "a1.jpg" )
	return ""
}

func fu2([]string) string {

	draw(acos)
	
	pImg := egui.GetWidg("main.img")
	pImg.SetImage( "a1.jpg" )
	return ""
}

func fu3([]string) string {

	draw(sqrt)
	
	pImg := egui.GetWidg("main.img")
	pImg.SetImage( "a1.jpg" )
	return ""
}

func fu4([]string) string {

	draw(newton)
	
	pImg := egui.GetWidg("main.img")
	pImg.SetImage( "a1.jpg" )
	return ""
}

func draw( fu func(z complex128) color.Color) {
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
