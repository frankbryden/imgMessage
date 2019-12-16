package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"strconv"
)

func getInts(msg string) []int {
	ints := make([]int, len(msg))
	for i, c := range "abc" {
		ints[i] = int(c)
	}
	return ints
}

func getBytes(ints []int) string {
	strReprBinaryVals := ""
	for _, val := range ints {
		strReprBinaryVals += strconv.FormatInt(int64(val), 2)
	}
	return strReprBinaryVals
}

func toggleOddEven(v *uint32) {
	if *v == 255 {
		*v--
	}
	*v++
}

func convertImage(img image.Image) *image.RGBA {
	b := img.Bounds()
	i := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(i, i.Bounds(), img, b.Min, draw.Src)
	return i
}

func processImage(img *image.RGBA, strBs string) {
	bounds := img.Bounds()

	x, y, i := 0, 0, 0

	for i < len(strBs) {
		r, g, b, _ := img.At(x, y).RGBA()
		r = r >> 8
		g = g >> 8
		b = b >> 8
		cols := []uint32{r, g, b}
		//we need to go through each channel: check if it corresponds with current bit
		//0 = even, 1 = odd
		fmt.Println("before", cols)
		for j, c := range cols {
			curVal, _ := strconv.Atoi(string(strBs[i]))
			fmt.Println("curVal", curVal)
			if int(c%2) != curVal {
				toggleOddEven(&c)
				cols[j] = c
			}

			i++
			if i >= len(strBs) {
				break
			}
		}
		x++
		if x == bounds.Max.X {
			x = 0
			y++
		}
		fmt.Println("after", cols)
		img.Set(x, y, color.RGBA{uint8(cols[0]), uint8(cols[1]), uint8(cols[2]), 0})
		//break
	}

	fmt.Println("w", bounds.Max.X, "h", bounds.Max.Y)
}

func main() {
	fmt.Println("hey")
	fileName := os.Args[1]
	message := os.Args[2]
	fmt.Println("Writing message", message, "to", fileName)

	reader, err := os.Open(fileName)

	if err != nil {
		fmt.Println("could not open file", err)
		os.Exit(1)
	}

	img, err := jpeg.Decode(reader)

	if err != nil {
		fmt.Println("failed to read file", err)
		os.Exit(1)
	}

	ints := getInts(message)

	//This is a list of string, where each string is a binary representation of each character in the message
	strBs := getBytes(ints)

	//Convert image.Image to image.RGBA
	rgba := convertImage(img)

	//Process img with bytes
	processImage(rgba, strBs)

	//
	f, _ := os.Create("out.jpg")
	e := jpeg.Encode(f, rgba, nil)
	
}
