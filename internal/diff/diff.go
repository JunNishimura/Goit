package diff

type Diff struct {
	textLong     []rune
	textShort    []rune
	lenTextLong  int
	lenTextShort int
	isReversed   bool
	EditDistance int
}

func NewDiff(text1, text2 []rune) *Diff {
	diff := &Diff{}
	if len(text1) > len(text2) {
		diff.textLong = text1
		diff.textShort = text2
		diff.lenTextLong = len(text1)
		diff.lenTextShort = len(text2)
	} else {
		diff.textLong = text2
		diff.textShort = text1
		diff.lenTextLong = len(text2)
		diff.lenTextShort = len(text1)
		diff.isReversed = true
	}

	return diff
}

func (d *Diff) Compose() {
	delta := d.lenTextLong - d.lenTextShort
	offset := d.lenTextShort + 1
	fp := make([]int, d.lenTextLong+d.lenTextShort+3)
	for i := range fp {
		fp[i] = -1
	}

	for p := 0; ; p++ {
		for k := -p; k < delta; k++ {
			fp[k+offset] = d.snake(k, max(fp[k+offset-1]+1, fp[k+offset+1]))
		}
		for k := delta + p; k > delta; k-- {
			fp[k+offset] = d.snake(k, max(fp[k+offset-1]+1, fp[k+offset+1]))
		}
		fp[delta+offset] = d.snake(delta, max(fp[delta+offset-1]+1, fp[delta+offset+1]))
		if fp[delta+offset] == d.lenTextLong {
			d.EditDistance = p
			break
		}
	}
}

func (d *Diff) snake(k, y int) int {
	x := y - k
	for x < d.lenTextShort && y < d.lenTextLong && d.textLong[y] == d.textShort[x] {
		x++
		y++
	}

	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
