package translation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Aligner struct {
	EnglishFile   string
	ForeignFile   string
	AlignmentFile string
	Tparams       map[string]map[string]float64
	Qparams       map[Q]float64
}

// Get tparams from file1, get qparams from fil2.
func (aligner *Aligner) GetParams(file1, file2 string) error {
	f1, err := os.Open(file1)
	if err != nil {
		return err
	}

	defer f1.Close()
	r1 := bufio.NewReader(f1)

	line, err := r1.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}
	err = json.Unmarshal([]byte(line), &aligner.Tparams)
	if err != nil {
		return err
	}

	f2, err := os.Open(file2)
	if err != nil {
		return err
	}

	defer f2.Close()
	r2 := bufio.NewReader(f2)
	aligner.Qparams = make(map[Q]float64)

	for {

		line, err := r2.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		args := strings.Split(strings.TrimSpace(line), " ")
		// i, j, l, m
		i, _ := strconv.Atoi(args[0])
		j, _ := strconv.Atoi(args[1])
		l, _ := strconv.Atoi(args[2])
		m, _ := strconv.Atoi(args[3])
		prob, _ := strconv.ParseFloat(args[4], 0)
		aligner.Qparams[Q{i: i, j: j, l: l, m: m}] = prob
	}
	return nil
}

// Compute best alightment given t params
func (aligner *Aligner) BestAlignment() error {
	fmt.Printf("Starting finding best alignment\n")
	ff, err := os.Open(aligner.ForeignFile)
	if err != nil {
		return err
	}

	ef, err := os.Open(aligner.EnglishFile)
	if err != nil {
		return err
	}

	af, err := os.Create(aligner.AlignmentFile)
	if err != nil {
		return err
	}

	defer ff.Close()
	defer ef.Close()
	defer af.Close()

	fr := bufio.NewReader(ff)
	er := bufio.NewReader(ef)

	k := 1
	for {
		engLine, err := er.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		forLine, err := fr.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		engWords := strings.Split(strings.TrimSpace(engLine), " ")
		forWords := strings.Split(strings.TrimSpace(forLine), " ")

		el := len(engWords)
		ml := len(forWords)
		i := 1
		for _, fw := range forWords {
			j := aligner.argmax2(i, el, ml, fw, &engWords)
			fmt.Fprintf(af, "%d %d %d\n", k, j, i)
			i++
		}
		k++
	}
	fmt.Printf("Ending finding best alignments\n")
	return nil
}
