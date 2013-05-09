<h1>simple-nlp</h1>

Simple NLP library, will include PCFG, Translation IBM 1 and 2 Models, GLM, etc...

Done through nlp class on coursera

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
    
Ex here soon...
    

