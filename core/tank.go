package tank

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
)

type grayMatrix = [][]uint8

func readToGray(filename string) (grayMatrix, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	gray := make(grayMatrix, height)
	for y := 0; y < height; y++ {
		gray[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			rgba := img.At(x, y)
			r, _, _, _ := color.GrayModel.Convert(rgba).RGBA()
			gray[y][x] = uint8(r >> 8)
		}
	}
	return gray, nil
}

func darken(img *grayMatrix, isTop bool) error {
	w, h := len((*img)[0]), len(*img)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			(*img)[y][x] >>= 1
			if isTop {
				(*img)[y][x] += 128
			}
		}
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sameSizeAndCenter(topImg, bottomImg grayMatrix) (grayMatrix, grayMatrix) {
	wT, hT := len(topImg[0]), len(topImg)
	wB, hB := len(bottomImg[0]), len(bottomImg)

	w, h := max(wT, wB), max(hT, hB)
	wm := abs(wT-wB) / 2
	hm := abs(hT-hB) / 2

	tRet := make(grayMatrix, h)
	bRet := make(grayMatrix, h)
	for y := 0; y < h; y++ {
		tRet[y] = make([]uint8, w)
		for x := 0; x < w; x++ {
			tRet[y][x] = 255
		}
		bRet[y] = make([]uint8, w)
	}

	// 假设 tRet 和 bRet 已经被初始化并且具有对应的切片类型
	// topImg 和 bottomImg 也是切片类型

	if w == wT && h == hT {
		copy(tRet[:hT], topImg) // 复制 top_img 到 t_ret
		for i := 0; i < hB; i++ {
			copy(bRet[hm+i][wm:wm+wB], bottomImg[i])
		}
	} else if w == wT && h == hB {
		for i := 0; i < hT; i++ {
			copy(tRet[hm+i][:wT], topImg[i])
		}
		copy(bRet[:hB], bottomImg) // 复制 bottom_img 到 b_ret
	} else if w == wB && h == hT {
		for i := 0; i < hT; i++ {
			copy(tRet[i][:wT], topImg[i])
		}
		for i := 0; i < hB; i++ {
			copy(bRet[hm+i][:wB], bottomImg[i])
		}
	} else {
		for i := 0; i < hT; i++ {
			for j := 0; j < wT; j++ {
				tRet[hm+i][wm+j] = topImg[i][j]
			}
		}
		for i := 0; i < hB; i++ {
			for j := 0; j < wB; j++ {
				bRet[i][j] = bottomImg[i][j]
			}
		}
	}

	return tRet, bRet
}

func merge(topImg, bottomImg grayMatrix) image.Image {
	w, h := len(topImg[0]), len(topImg)
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			alpha := 255 - (topImg[y][x] - bottomImg[y][x])
			var gray uint8
			if alpha == 255 {
				alpha = 0
			} else {
				gray = uint8(float32(bottomImg[y][x]) / float32(alpha) * 255)
			}
			c := color.NRGBA{R: gray, G: gray, B: gray, A: alpha}
			newImg.Set(x, y, c)
		}
	}

	return newImg
}

func Make(topImg, bottomImg, output string) (err error) {
	tg, err := readToGray(topImg)
	if err != nil {
		return
	}
	if darken(&tg, true) != nil {
		return
	}
	bg, err := readToGray(bottomImg)
	if err != nil {
		return
	}
	if darken(&bg, false) != nil {
		return
	}

	tg2, bg2 := sameSizeAndCenter(tg, bg)

	img := merge(tg2, bg2)
	f, err := os.Create(output)
	if err != nil {
		return
	}
	defer f.Close()
	png.Encode(f, img)

	return
}
