package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/po3rin/gonbayes"
)

const (
	posiLabel = "positive"
	negaLabel = "negative"
)

func loadDataset(file string) (map[string]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dataset := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		data := strings.Split(l, "\t")
		if len(data) != 2 {
			continue
		}
		s := data[0]
		if data[1] == "0" {
			dataset[s] = negaLabel
		} else if data[1] == "1" {
			dataset[s] = posiLabel
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return dataset, nil
}

func main() {
	s := flag.String("s", "", "input string")
	f := flag.String("f", "./yelp_labelled.txt", "dataset file path")
	t := flag.String("t", "", "trained model file path")
	o := flag.String("o", "negaposi_classifier.gob", "file name for output trained model")
	flag.Parse()

	class := []string{posiLabel, negaLabel}
	classifier := gonbayes.NewClassifier(class)

	if *t != "" {
		err := classifier.Decode(*t)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		dataset, err := loadDataset(*f)
		if err != nil {
			log.Fatal(err)
		}

		for s, v := range dataset {
			classifier.Train(v, s)
		}
	}

	result := classifier.Classify(*s)
	fmt.Println(result)

	if *t == "" {
		classifier.Encode(*o)
	}
}
