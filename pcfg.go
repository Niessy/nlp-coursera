package pcfg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type PCFG struct {
	wordCounts        map[string]float64
	nonTerminalCounts map[NonTerminal]float64
	binaryCounts      map[BinaryRule]float64
	unaryCounts       map[UnaryRule]float64
	rareThreshold     int
}

func NewPCFG(threshold int) *PCFG {
	p := new(PCFG)
	pp := *p
	pp.wordCounts = make(map[string]float64)
	pp.nonTerminalCounts = make(map[NonTerminal]float64)
	pp.binaryCounts = make(map[BinaryRule]float64)
	pp.unaryCounts = make(map[UnaryRule]float64)
	pp.rareThreshold = threshold
	*p = pp
	return p
}

type BinaryRule struct {
	X      string
	Y1, Y2 string
}

type UnaryRule struct {
	X, Y string
}

type NonTerminal struct {
	NT string
}

type SubTree struct {
	Start int
	End   int
	Root  string
}

type BPSubTree struct {
	Root  string
	Left  SubTree
	Right SubTree
}

// Takes a file with terminal counts and stores them in
// PCFG.
func (pcfg *PCFG) GetWordCounts(countfile string) {
	f0, err := os.Open(countfile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f0.Close()

	cr := bufio.NewReader(f0)
	for {
		line, err := cr.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.Replace(line, "\n", "", -1)
		s := strings.Split(line, " ")
		count, _ := strconv.ParseFloat(s[0], 0)
		switch s[1] {
		case "UNARYRULE":
			pcfg.incCount(s[3], count)
			unary := UnaryRule{s[2], s[3]}
			pcfg.incCount(unary, count)
			continue
		case "BINARYRULE":
			binary := BinaryRule{s[2], s[3], s[4]}
			pcfg.incCount(binary, count)
			continue
		case "NONTERMINAL":
			nt := NonTerminal{s[2]}
			pcfg.incCount(nt, count)
			continue
		}
	}
}

// Main routine
func (pcfg *PCFG) RewriteTrainingTree(trainingfile, resultfile string) error {
	f1, err := os.Open(trainingfile)
	if err != nil {
		return err
	}
	defer f1.Close()

	tc := bufio.NewReader(f1)

	f2, err := os.Create(resultfile)
	if err != nil {
		return err
	}
	defer f2.Close()

	// Run CKY on the training file
	counter := 1
	for {
		line, err := tc.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1)
		s := strings.Split(line, " ")
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Printf("Running CKY on sentence %d.\n", counter)
		tree := pcfg.ckyAlgorithm(s, "SBARQ")
		b, err := json.Marshal(&tree)
		if err != nil {
			return err
		}
		fmt.Printf("Writing result for sentence %d\n to file.", counter)
		fmt.Fprintf(f2, "%s\n", b)
		counter++
	}
	return nil
}

func (pcfg *PCFG) incCount(key interface{}, num float64) {
	switch vv := key.(type) {
	case string:
		if _, ok := pcfg.wordCounts[vv]; !ok {
			pcfg.wordCounts[vv] = num
		} else {
			pcfg.wordCounts[vv] += num
		}

	case UnaryRule:
		if _, ok := pcfg.unaryCounts[vv]; !ok {
			pcfg.unaryCounts[vv] = num
		} else {
			pcfg.unaryCounts[vv] += num
		}

	case BinaryRule:
		if _, ok := pcfg.binaryCounts[vv]; !ok {
			pcfg.binaryCounts[vv] = num
		} else {
			pcfg.binaryCounts[vv] += num
		}

	case NonTerminal:
		if _, ok := pcfg.nonTerminalCounts[vv]; !ok {
			pcfg.nonTerminalCounts[vv] = num
		} else {
			pcfg.nonTerminalCounts[vv] += num
		}
	}
}

// ckyAlgorithm
func (pcfg *PCFG) ckyAlgorithm(sentence []string, root string) []interface{} {
	pi := make(map[SubTree]float64)
	bp := make(map[SubTree]interface{})

	// Initialization
	for i, word := range sentence {
		rare := true
		for nt, _ := range pcfg.nonTerminalCounts {
			st := SubTree{i + 1, i + 1, nt.NT}
			currentRule := UnaryRule{nt.NT, word}
			if _, ok := pcfg.unaryCounts[currentRule]; ok {
				pi[st] = pcfg.q(currentRule)
				bp[st] = currentRule
				rare = false
			}
		}
		if rare {
			for nt, _ := range pcfg.nonTerminalCounts {
				st := SubTree{i + 1, i + 1, nt.NT}
				rareRule := UnaryRule{nt.NT, "_RARE_"}
				pi[st] = pcfg.q(rareRule)
				bp[st] = UnaryRule{nt.NT, word}
			}
		}
	}

	// Main part of the algorithm
	n := len(sentence)
	for l := 1; l < n; l++ {
		for i := 1; i <= n-l; i++ {
			j := i + l
			for nt, _ := range pcfg.nonTerminalCounts {
				st := SubTree{i, j, nt.NT}
				max_prob, max_arg := pcfg.maxSubTree(nt.NT, i, j, pi)
				pi[st] = max_prob
				bp[st] = max_arg
			}
		}
	}

	// Return the max prob tree
	_, mta := pcfg.maxSubTree(root, 1, n, pi)
	var tree []interface{}
	finalTree := getTree(mta, tree, bp)
	return finalTree
}

// Recursively fill the PCFG based on max arg backpointers.
func getTree(st interface{}, tree []interface{}, bp map[SubTree]interface{}) []interface{} {
	switch vv := st.(type) {
	case UnaryRule:
		tree = append(tree, vv.X)
		tree = append(tree, vv.Y)
	case BPSubTree:
		left := getTree(bp[vv.Left], tree, bp)
		right := getTree(bp[vv.Right], tree, bp)
		tree = append(tree, vv.Root)
		tree = append(tree, left)
		tree = append(tree, right)
	}
	return tree
}

func (pcfg *PCFG) maxSubTree(x string, i, j int, pi map[SubTree]float64) (max float64, arg BPSubTree) {
	for br, _ := range pcfg.binaryCounts {
		// We have a match
		if br.X == x {
			derivProb := pcfg.q(br)
			for s := i; s < j; s++ {
				left := pi[SubTree{i, s, br.Y1}]
				right := pi[SubTree{s + 1, j, br.Y2}]
				p := derivProb * left * right
				if p > max {
					max = p
					arg = BPSubTree{br.X, SubTree{i, s, br.Y1}, SubTree{s + 1, j, br.Y2}}
				}
			}
		}
	}
	return
}

func (pcfg *PCFG) q(rule interface{}) (f float64) {
	var nt NonTerminal
	switch vv := rule.(type) {
	case BinaryRule:
		nt = NonTerminal{vv.X}
		f = pcfg.binaryCounts[vv] / pcfg.nonTerminalCounts[nt]
		break
	case UnaryRule:
		nt = NonTerminal{vv.X}
		f = pcfg.unaryCounts[vv] / pcfg.nonTerminalCounts[nt]
	}
	return
}
