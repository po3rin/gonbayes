// Package gonbayes is Simple Naive Bayes Classifier in Go.
package gonbayes

import (
	"encoding/gob"
	"log"
	"os"
	"sort"

	"github.com/pkg/errors"
)

// Classifier is documents types categoclasifier.
type Classifier struct {
	Words             map[string]map[string]uint64
	TotalWords        uint64
	TotalDocsInTypes  map[string]uint64
	TotalDocs         uint64
	TotalWordsInTypes map[string]uint64
	Threshold         float64
}

// NewClassifier inits classifier.
func NewClassifier(types []string, Threshold float64) (c Classifier) {
	c = Classifier{
		Words:             make(map[string]map[string]uint64),
		TotalWords:        0,
		TotalDocsInTypes:  make(map[string]uint64),
		TotalDocs:         0,
		TotalWordsInTypes: make(map[string]uint64),
		Threshold:         Threshold,
	}

	for _, category := range types {
		c.Words[category] = make(map[string]uint64)
		c.TotalDocsInTypes[category] = 0
		c.TotalWordsInTypes[category] = 0
	}
	return
}

// Train trains documents classifier.
func (c *Classifier) Train(category string, document string) {
	for word, count := range countWords(document) {
		c.Words[category][word] += uint64(count)
		c.TotalWordsInTypes[category] += uint64(count)
		c.TotalWords += uint64(count)
	}
	c.TotalDocsInTypes[category]++
	c.TotalDocs++
}

func (c *Classifier) pCategory(category string) float64 {
	return float64(c.TotalDocsInTypes[category]) / float64(c.TotalDocs)
}

func (c *Classifier) pDocCategory(category string, document string) (p float64) {
	p = 1.0
	for word := range countWords(document) {
		p = p * c.pWordCategory(category, word)
	}
	return p
}

func (c *Classifier) pWordCategory(category string, word string) float64 {
	return float64(c.Words[category][stem(word)]+1) / float64(c.TotalWordsInTypes[category])
}

func (c *Classifier) pCategoryDocument(category string, document string) float64 {
	return c.pDocCategory(category, document) * c.pCategory(category)
}

// P is Probabilities of each types.
func (c *Classifier) P(document string) (p map[string]float64) {
	p = make(map[string]float64)
	for category := range c.Words {
		p[category] = c.pCategoryDocument(category, document)
	}
	return
}

type sorted struct {
	category    string
	probability float64
}

// Classify classify documents.
func (c *Classifier) Classify(document string) (category string) {
	prob := c.P(document)

	var sp []sorted
	for c, p := range prob {
		sp = append(sp, sorted{c, p})
	}
	sort.Slice(sp, func(i, j int) bool {
		return sp[i].probability > sp[j].probability
	})

	if sp[0].probability/sp[1].probability > c.Threshold {
		category = sp[0].category
	} else {
		category = "unknown"
	}

	return
}

// Encode trained Classifier
func (c *Classifier) Encode(fileName string) error {
	if c.TotalDocs == 0 {
		return errors.New("gonbayes: classifier is not trained yet")
	}
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = gob.NewEncoder(f).Encode(&c)
	if err != nil {
		return errors.Wrap(err, "gonbayes: failed to encode")
	}
	return nil
}

// Decode CBOW output file to struct.
func (c *Classifier) Decode(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = gob.NewDecoder(f).Decode(c)
	if err != nil {
		return errors.Wrap(err, "gonbayes: failed to dencode file")
	}
	return nil
}
