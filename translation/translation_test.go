package translation

import (
	"testing"
)

var err error

// func TestAligner(t *testing.T) {
// 	a := new(Aligner)
// 	a.Tparams = make(map[string]map[string]float64)
// 	a.EnglishFile = "../../dev.en"
// 	a.ForeignFile = "../../dev.es"
// 	a.AlignmentFile = "dev.p1.out"
// 	err = a.GetParams("dev.json2")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = a.BestAlignment()
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestIBM2Init(t *testing.T) {
	i := new(IBM2)
	i.Ibm1.ForeignFile = "../../corpus.es"
	i.Ibm1.EnglishFile = "../../corpus.en"
	i.TparamInput = "dev.json"
	i.QparamsFile = "q1.json"
	i.Ibm1.TparamsFile = "t1.json"
	i.Qparams = make(map[Q]float64)
	err := i.Initialize()
	if err != nil {
		t.Error(err)
	}
	err = i.EMAlgorithm(5)
	if err != nil {
		t.Error(err)
	}
}
