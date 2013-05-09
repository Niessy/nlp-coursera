package translation

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func NewIBM1(ef, ff, tpf string) *IBM1 {
	m := new(IBM1)
	mm := *m
	mm.EnglishFile = ef
	mm.ForeignFile = ff
	mm.TparamsFile = tpf
	mm.Tparams = make(map[string]map[string]float64)
	*m = mm
	return m
}

type IBM1 struct {
	EnglishFile string
	ForeignFile string
	TparamsFile string
	Tparams     map[string]map[string]float64
}

func (model *IBM1) Initialize() error {

	ff, err := os.Open(model.ForeignFile)
	if err != nil {
		return err
	}

	ef, err := os.Open(model.EnglishFile)
	if err != nil {
		return err
	}

	defer ff.Close()
	defer ef.Close()

	fr := bufio.NewReader(ff)
	er := bufio.NewReader(ef)

	model.Tparams["NULL"] = make(map[string]float64)

	// Set t(f | e) = 1 / n(e)
	for {
		engLine, err := er.ReadString('\n')
		if checkerror(err) {
			normalize(model)
			return err
		}

		forLine, err := fr.ReadString('\n')
		if checkerror(err) {
			normalize(model)
			return err
		}

		if engLine == "\n" || forLine == "\n" {
			continue
		}

		engWords := strings.Split(strings.TrimSpace(engLine), " ")
		forWords := strings.Split(strings.TrimSpace(forLine), " ")

		for _, ew := range engWords {
			for _, fw := range forWords {
				if _, ok := model.Tparams[ew]; !ok {
					model.Tparams[ew] = make(map[string]float64)
					model.Tparams[ew][fw] = 1
				} else if _, ok := model.Tparams[ew][fw]; !ok {
					model.Tparams[ew][fw] = 1
				}

				if _, ok := model.Tparams["NULL"][fw]; !ok {
					model.Tparams["NULL"][fw] = 1
				}
			}
		}

	}
	return nil
}

// Runs EM Algorithm for n iterations then
// writes tparams to TparamsFile
func (model *IBM1) EMAlgorithm(n int) error {
	tt := time.Now()
	fmt.Printf("Beginning EM for IBM1\n")
	for s := 1; s <= n; s++ {
		fmt.Printf("Beginning iteration %d...\n", s)
		t1 := time.Now()
		// Reset counts
		counts := make(map[string]float64)

		// being reading files
		ff, err := os.Open(model.ForeignFile)
		if err != nil {
			return err
		}

		ef, err := os.Open(model.EnglishFile)
		if err != nil {
			return err
		}

		er := bufio.NewReader(ef)
		fr := bufio.NewReader(ff)

		for {
			engLine, err := er.ReadString('\n')
			if checkerror(err) {
				break
			}

			forLine, err := fr.ReadString('\n')
			if checkerror(err) {
				break
			}

			if engLine == "\n" || forLine == "\n" {
				continue
			}

			engWords := strings.Split("NULL "+strings.TrimSpace(engLine), " ")
			forWords := strings.Split(strings.TrimSpace(forLine), " ")

			i := 1
			for _, fw := range forWords {
				for _, ew := range engWords {
					d := model.delta(fw, ew, &engWords)
					counts[ew+" "+fw] = counts[ew+" "+fw] + d
					counts[ew] = counts[ew] + d
				}
				i++
			}

		}
		// Cleanup files for next iteration
		ff.Close()
		ef.Close()

		// Revising Tparams
		for e, fws := range model.Tparams {
			for f, _ := range fws {
				model.Tparams[e][f] = counts[e+" "+f] / counts[e]
			}
		}

		fmt.Printf("Ending iteration %d...took %v\n", s, time.Since(t1))
	}
	fmt.Println("Finished EM Algorithm took ", time.Since(tt))
	fmt.Printf("\nWriting tparams to %s...\n", model.TparamsFile)
	err := model.tparamsToFile()
	if err != nil {
		return err
	}
	fmt.Printf("Finished writing tparams...\n")
	return nil
}
