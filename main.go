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

//Flip parity
func toggleOddEven(v *uint32) {
	if *v == 255 {
		*v--
	} else {
		*v++
	}
}

func convertImage(img image.Image) *image.RGBA {
	b := img.Bounds()
	i := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(i, i.Bounds(), img, b.Min, draw.Src)
	return i
}

func decodeImage(img *image.RGBA) {

	//Create a string of 0's and 1's that has the length of the message
	messageLength := ""
	x, y := 0, 0
	bounds := img.Bounds()

	//Get the values at the first 12 channels (4 pixels)
	for x < 4 {
		r, g, b, _ := img.At(x, y).RGBA()
		r = r >> 8
		g = g >> 8
		b = b >> 8
		messageLength += string(r%2) + string(g%2) + string(b%2)
		x++
	}

	//Parse the message length string to an integer
	length, _ := strconv.ParseInt(messageLength, 2, 32)

	bitMessage := ""
	bitCounter := 0

	for bitCounter < int(length) {
		r, g, b, _ := img.At(x, y).RGBA()
		r = r >> 8
		g = g >> 8
		b = b >> 8
		cols := []uint32{r, g, b}
		//we need to go through each channel: check if it corresponds with current bit
		//0 = even, 1 = odd
		for _, c := range cols {
			bitMessage += string(c % 2)
			bitCounter++
			if bitCounter >= len(messageLength) {
				break
			}
		}
		x++
		if x >= bounds.Max.X {
			x = 0
			y++
		}
	}

}

//Takes image and message to write into image
func processImage(img *image.RGBA, strBs string) {
	bounds := img.Bounds()

	//i = current index in the message
	x, y, i := 0, 0, 0

	for i < len(strBs) {
		//Get the RGB values at (x,y)
		r, g, b, _ := img.At(x, y).RGBA()
		//Make the pixel values between 0 and 255
		r = r >> 8
		g = g >> 8
		b = b >> 8
		cols := []uint32{r, g, b}
		//we need to go through each channel: check if it corresponds with current bit
		//0 = even, 1 = odd
		for j, c := range cols {
			//Convert character to string, then to an integer (0 or 1)
			curVal, _ := strconv.Atoi(string(strBs[i]))
			fmt.Println("curVal", curVal)
			//Check parity of channel against parity of current bit in message
			if int(c%2) != curVal {
				toggleOddEven(&c)
				cols[j] = c
			}

			i++
			if i >= len(strBs) {
				break
			}
		}
		img.Set(x, y, color.RGBA{uint8(cols[0]), uint8(cols[1]), uint8(cols[2]), 0})
		x++
		if x == bounds.Max.X {
			x = 0
			y++
		}
	}
}

func main() {
	fmt.Println("hey")
	fileName := os.Args[1]
	message := ""
	encode := false
	if len(os.Args) > 2 {
		encode = true
		message = os.Args[2]
		fmt.Println("Writing message", message, "to", fileName)
	}

	reader, err := os.Open(fileName)

	if err != nil {
		fmt.Println("could not open file", err)
		os.Exit(1)
	}

	//Read the image
	img, err := jpeg.Decode(reader)

	if err != nil {
		fmt.Println("failed to read file", err)
		os.Exit(1)
	}

	//Here, we do 2 different things, based on the value of encode: call Frankie's code if false, and Jasmine's code if true

	//Get the ASCII value for each character in the message
	ints := getInts(message)

	//This is a list of string, where each string is a binary representation of each character in the message
	//Convert each ASCII value into binary
	//Concatenate everything into one whole string
	strBs := getBytes(ints)

	//Convert image.Image (read-only) to image.RGBA (editable format)
	rgba := convertImage(img)

	//Process img with bytes
	processImage(rgba, strBs)

	//
	f, _ := os.Create("out.jpg")
	e := jpeg.Encode(f, rgba, nil)
	fmt.Println("Success", e)
}
