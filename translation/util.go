package translation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Transfer Tparams to file in json format for
// later use in an Aligner.
func (model *IBM1) tparamsToFile() error {
	f, err := os.Create(model.TparamsFile)
	if err != nil {
		return err
	}
	b, err := json.Marshal(model.Tparams)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "%s", b)
	return nil
}

// Computes best alignment based on previous IBM 1 models.
func (model *IBM1) delta(fw, ew string, engWords *[]string) float64 {
	bot := 0.0
	for _, w := range *engWords {
		bot += model.Tparams[w][fw]
	}
	return model.Tparams[ew][fw] / bot
}

// Computes best alignment based on previous IBM 2 models.
func (model *IBM2) delta2(i, k, l, m int, ew, fw string, engWords *[]string) float64 {
	bot := 0.0
	for j, w := range *engWords {
		q := Q{i: i, j: j, l: l, m: m}
		bot += (model.Ibm1.Tparams[w][fw] * model.Qparams[q])
	}
	q := Q{i: i, j: k, l: l, m: m}
	return (model.Ibm1.Tparams[ew][fw] * model.Qparams[q]) / bot
}

// Transform n(e) to 1 / n(e)
func normalize(model *IBM1) {
	for _, keys := range model.Tparams {
		n := (1 / float64(len(keys)))
		for key, _ := range keys {
			keys[key] = n
		}
	}
}

// argmax for IBM Model 2
func (aligner *Aligner) argmax2(i, l, m int, fw string, engWords *[]string) (n int) {
	max := 0.0
	for j, ew := range *engWords {
		q := Q{i: i, j: j + 1, l: l, m: m}
		prob := aligner.Tparams[ew][fw] * aligner.Qparams[q]
		if prob > max {
			n = j + 1
			max = prob
		}
	}
	return
}

func (model *IBM2) readtparams() error {
	f, err := os.Open(model.TparamInput)
	if err != nil {
		return err
	}

	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}
	err = json.Unmarshal([]byte(line), &model.Ibm1.Tparams)
	if err != nil {
		return err
	}
	return nil
}

func (model *IBM2) qparamsToFile() error {
	f, err := os.Create(model.QparamsFile)
	if err != nil {
		return err
	}
	for qval, prob := range model.Qparams {
		fmt.Fprintf(f, "%d %d %d %d %f\n", qval.i, qval.j, qval.l, qval.m, prob)
	}
	return nil
}
