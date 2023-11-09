package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Configuration struct {
	Corpuspath []string `json:"corpuspath"`
}

var (
	Conf Configuration
)

func init() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		panic("Unable to load config.json :" + err.Error())
	}
	defer jsonFile.Close()

	byteData, err := io.ReadAll(jsonFile)
	if err != nil {
		panic("Unable to read jsonFile :" + err.Error())
	}

	err = json.Unmarshal(byteData, &Conf)
	if err != nil {
		panic("Unable to unmarshall jsonFile -> byteData :" + err.Error())
	}
}

func main() {
	/* init map */
	Emit = make(map[string]float64, 0)
	Transition = make(map[string]float64, 0)
	Context = make(map[string]float64, 0)

	model := MarkovModel{
		Emit:           Emit,
		Transition:     Transition,
		Context:        Context,
		TransitionData: make([]MarkovTransition, 0),
	}

	for _, path := range Conf.Corpuspath {

		corpora, err := GetCorpora(path)
		if err != nil {
			err = errors.New(fmt.Sprintf("corpora: %v is error. %v", path, err.Error()))
			// skip corpora
			fmt.Println(err.Error())
		}
		defer corpora.Close()

		/* open new scanner for reading each line  */
		scanner := bufio.NewScanner(corpora)
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

						model.Transition[previous+" "+tag]++
						model.Context[tag]++
						model.Emit[tag+" "+word]++
						previous = tag
					}
				}
				model.Transition[previous+"</s>"]++
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
			td := MarkovTransition{
				POSKey:        key,
				PreviousPOS:   previous,
				POS:           word,
				Probabilities: TransitionProbabilities,
			}

			model.TransitionData = append(model.TransitionData, td)

			fmt.Printf("\n")
		}

		/* print the transition probabilities*/
		for key, value := range Emit {
			//fmt.Printf("key,value : %v,%v", key, value)

			previous_tag = strings.Split(key, " ")
			if len(previous_tag) != 2 {
				// Handle the case where the key does not contain a space
				continue
			}

			tag = previous_tag[0]
			word = previous_tag[1]

			// if the tag is A , there a % that the word will e X

			// Print the information
			if Context[tag] == 0 {
				EmissionProbabilities = 0
				fmt.Printf("T %v:%v", key, EmissionProbabilities)
			} else {
				EmissionProbabilities = value / Context[tag]
				fmt.Printf("T %v:%v", key, EmissionProbabilities)
			}
			td := MarkovEmission{
				POSKey:        key,
				Tag:           tag,
				Word:          word,
				Probabilities: EmissionProbabilities,
			}

			model.EmissionData = append(model.EmissionData, td)

			fmt.Printf("\n")
		}
	}

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(model.TransitionData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Marshal the struct to JSON
	jsonEmission, err := json.Marshal(model.EmissionData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	saveJsonToFile("transition_probabilities", jsonData)
	saveJsonToFile("emission_probabilities", jsonEmission)
}
