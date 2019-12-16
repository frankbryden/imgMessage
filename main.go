package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
)

func padBitString(b string, n int) string {
	for len(b) < n {
		b = "0" + b
	}
	return b
}

func getInts(msg string) []int {
	ints := make([]int, len(msg))
	for i, c := range msg {
		ints[i] = int(c)
	}
	return ints
}

func getBytes(ints []int) string {
	strReprBinaryVals := ""
	for _, val := range ints {
		s := strconv.FormatInt(int64(val), 2)
		s = padBitString(s, 7)
		strReprBinaryVals += s
	}
	return strReprBinaryVals
}

//read color.Color and return bit string of length 3 (RGB)
func readPixel(col color.Color) string {
	r, g, b, _ := col.RGBA()
	r = r >> 8
	g = g >> 8
	b = b >> 8
	cols := []uint32{r, g, b}
	bitString := ""
	for _, c := range cols {
		if c%2 == 0 {
			bitString += "0"
		} else {
			bitString += "1"
		}
	}
	return bitString
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

func decodeMessage(img *image.RGBA) {
	//WHAT REALLY NEEDS TO BE DONE: write a function which simply reads n number of bits from the image and returns them as a string
	//Create a string of 0's and 1's that has the length of the message
	messageLength := ""
	x, y := 0, 0
	bounds := img.Bounds()

	//Get the values at the first 12 channels (4 pixels)
	for x < 4 {
		messageLength += readPixel(img.At(x, y))
		x++
	}

	//Parse the message length string to an integer
	l, _ := strconv.ParseInt(messageLength, 2, 32)
	length := int(l)

	fmt.Println("message length:", length)

	bitMessage := ""
	bitCounter := 0

	for bitCounter < length {
		messagePart := readPixel(img.At(x, y))
		if (length - bitCounter) <= 3 { //if we have less than 3 bits left to read, don't add entire message
			messagePart = messagePart[0 : length-bitCounter]
		}
		bitMessage += messagePart
		bitCounter += 3
		x++
		if x >= bounds.Max.X {
			x = 0
			y++
		}
	}
	//We now have a list of bits in the form of a string
	//looking a bit like this: '101100110101'
	//we need to split it in 7s, where each group of seven represents an int, in turn
	//representing an ASCII character
	if len(bitMessage)%7 != 0 {
		fmt.Println("error, bitMessage does not contain a multiple of 7 bits")
	}

	ints := make([]int, len(bitMessage)/7)
	chars := make([]rune, len(bitMessage)/7)
	for i := 0; i < len(bitMessage); i += 7 {
		part := bitMessage[i : i+7]
		asciiInt, _ := strconv.ParseInt(part, 2, 32)
		ints[i/7] = int(asciiInt)
		chars[i/7] = rune(asciiInt)
	}
	fmt.Println(string(chars))

}

func encodeMessage(message string, rgba *image.RGBA) {
	//Get the ASCII value for each character in the message
	ints := getInts(message)

	//This is a list of string, where each string is a binary representation of each character in the message
	//Convert each ASCII value into binary
	//Concatenate everything into one whole string
	strBs := getBytes(ints)

	//Process img with bytes
	processImage(rgba, strBs)

	//
	f, _ := os.Create("out.png")
	defer f.Close()
	e := png.Encode(f, rgba)
	fmt.Println("Success", e)
}

//Takes image and message to write into image
func processImage(img *image.RGBA, strBs string) {
	bounds := img.Bounds()

	//i = current index in the message
	x, y, i := 0, 0, 0

	//bitLength is simply the number of bits used to describe the message
	bitLength := getBytes([]int{len(strBs)})

	//make sure it is 12 bits long
	bitLength = padBitString(bitLength, 12)

	for j := 0; j < 4; j++ {
		r, g, b, a := img.At(x, y).RGBA()

		//Make the pixel values between 0 and 255
		r = r >> 8
		g = g >> 8
		b = b >> 8
		cols := []uint32{r, g, b}
		//we need to go through each channel: check if it corresponds with current bit
		//0 = even, 1 = odd
		for k, c := range cols {
			//Convert character to string, then to an integer (0 or 1)
			curVal, _ := strconv.Atoi(string(bitLength[3*j+k]))
			//Check parity of channel against parity of current bit in message
			if int(c%2) != curVal {
				toggleOddEven(&c)
				cols[k] = c
			}
		}
		img.Set(x, y, color.RGBA{uint8(cols[0]), uint8(cols[1]), uint8(cols[2]), uint8(a)})
		x++
		if x == bounds.Max.X {
			x = 0
			y++
		}

	}

	//write bitLength (number of bits composing message) to signify decoder how much to read
	fmt.Println("bitLength", bitLength)

	for i < len(strBs) {
		//Get the RGB values at (x,y)
		r, g, b, a := img.At(x, y).RGBA()
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
		img.Set(x, y, color.RGBA{uint8(cols[0]), uint8(cols[1]), uint8(cols[2]), uint8(a)})
		x++
		if x == bounds.Max.X {
			x = 0
			y++
		}
	}
}

func main() {
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
	img, err := png.Decode(reader)

	if err != nil {
		fmt.Println("failed to read file", err)
		os.Exit(1)
	}

	//Convert image.Image (read-only) to image.RGBA (editable format)
	rgba := convertImage(img)

	//Here, we do 2 different things, based on the value of encode: call Frankie's code if false, and Jasmine's code if true
	if encode {
		encodeMessage(message, rgba)
	} else {
		decodeMessage(rgba)
	}

}
