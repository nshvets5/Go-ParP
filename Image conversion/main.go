package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"sync"
	"time"
)

func convertSequential(inputImage image.Image) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	grayImage := image.NewGray(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			originalColor := inputImage.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImage.Set(x, y, grayColor)
		}
	}
	return grayImage
}

func convertParallel(inputImage image.Image) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	grayImage := image.NewGray(bounds)

	var wg sync.WaitGroup
	for y := 0; y < height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < width; x++ {
				originalColor := inputImage.At(x, y)
				grayColor := color.GrayModel.Convert(originalColor)
				grayImage.Set(x, y, grayColor)
			}
		}(y)
	}
	wg.Wait()

	return grayImage
}

func saveImage(filename string, img image.Image) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, img, nil)
}

func main() {
	inputFile, err := os.Open("input.jpg")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	start := time.Now()
	seqImg := convertSequential(img)
	err = saveImage("output_sequential.jpg", seqImg)
	if err != nil {
		fmt.Println("Error saving sequential image:", err)
		return
	}
	fmt.Printf("Sequential conversion took %v\n", time.Since(start))

	start = time.Now()
	parImg := convertParallel(img)
	err = saveImage("output_parallel.jpg", parImg)
	if err != nil {
		fmt.Println("Error saving parallel image:", err)
		return
	}
	fmt.Printf("Parallel conversion took %v\n", time.Since(start))
}
