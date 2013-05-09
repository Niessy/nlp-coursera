<h1>simple-nlp</h1>

Simple NLP library, will include PCFG, Translation IBM 1 and 2 Models, GLM, etc...

Done through nlp class on coursera

To use in your GOPATH, should be (YOUR GO DIR)/src/github.com/Niessy.
	
    git clone https://github.com/Niessy/simple-nlp.git
    
Or

    go get github.com/Niessy/simple-nlp

<h2> Probabilistic Context Free Grammar </h2>

Example of usage: This works assuming the counts file already has rare counts, I'm going to add support for this soon.

    package main

    import (
        "fmt"
	    "github.com/Niessy/simple-nlp/pcfg"
    )

    func main() {
	    p := pcfg.NewPCFG(5)
	    p.GetCounts("counts.test")
	    err := p.ParseSentences("parse_dev.dat", "output.json")
	    if err != nil {
		    fmt.Println(err)
	    }
    }
    
<h2> Machine Translation </h2>
    
Example:

	package main
	
	import (
		"github.com/Niessy/simple-nlp/translation"
		"fmt"
	)
	
	func main() {
		// corpus.en is the native text, corpus.es acts as the foreign text
		// output.json contains associated probabilities for aligned native/foreign words
		// used as input for IBM2
		i := translation.NewIBM1("corpus.en", "corpus.es", "output.json")
		err := i.Initialize()
		if err != nil {
			fmt.Println(err)
		}
		err = i.EMAlgorithm(5) // Number of iterations is 5
		if err != nil {
			fmt.Println(err)
		}
	}

    

