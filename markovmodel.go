package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	Emit       map[string]float64
	Transition map[string]float64
	Context    map[string]float64

	line          string
	wordtags      []string
	wordtag_split []string

	word string
	tag  string

	previous string

	previous_word []string

	TransitionProbabilities float64
	EmissionProbabilities   float64
)

type MarkovTransition struct {
	POSKey        string  `json:"poskey"`
	PreviousPOS   string  `json:"previous_pos"`
	POS           string  `json:"pos"`
	Probabilities float64 `json:"probabilities"`
}

type MarkovEmission struct {
}

type MarkovModel struct {
	Corpora        *os.File
	Emit           map[string]float64
	Transition     map[string]float64
	Context        map[string]float64
	TransitionData []MarkovTransition
}

func GetCorpora(path string) (corpora *os.File, err error) {
	/* init corpora (single file)  */
	// filename := `/home/jeremykenn/Documents/Kuliah S2/NLP/brown/ca01`
	// filename := `./brown/ca01`
	corpora, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	return corpora, nil
}

func markovmodel() {

	/* open new scanner for reading each line  */
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = scanner.Text()

		/* check if line is not null  */
		if len(line) > 0 { // berarti ada kalimat
			/* start sentence*/
			previous = "<s>"
			// remove first tab
			line = strings.Replace(line, "\t", "", -1)
			//fmt.Printf("%v. %v. %v\n", i, len(line), line)

			// split line into word tags
			wordtags := strings.Split(line, " ")

			for _, wordtag := range wordtags {
				//fmt.Printf("%v ", wordtag)
				// split wordtag into word and tag
				wordtag_split = strings.Split(wordtag, "/")

				//fmt.Printf("%v", len(wordtag_split))

				//fmt.Printf("%v", wordtag)
				//fmt.Printf("%v", wordtag_split)
				if len(wordtag_split) > 1 {

					word = wordtag_split[0]
					tag = wordtag_split[1]
					fmt.Printf("word:%v, tag:%v.", word, tag)

					Transition[previous+" "+tag]++
					Context[tag]++
					Emit[tag+" "+word]++
					previous = tag
				}
			}
			Transition[previous+"</s>"]++
			fmt.Printf("\n")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	/* print the transition probabilities*/
	for key, value := range Transition {
		//fmt.Printf("key,value : %v,%v", key, value)

		previous_word = strings.Split(key, " ")
		if len(previous_word) != 2 {
			// Handle the case where the key does not contain a space
			continue
		}

		previous = previous_word[0]
		word = previous_word[1]

		// Print the information
		if Context[previous] == 0 {
			TransitionProbabilities = 0
			fmt.Printf("T %v:%v", key, TransitionProbabilities)
		} else {
			TransitionProbabilities = value / Context[previous]
			fmt.Printf("T %v:%v", key, TransitionProbabilities)
		}

		fmt.Printf("\n")
	}
}
