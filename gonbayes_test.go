package gonbayes_test

import (
	"reflect"
	"testing"

	"github.com/po3rin/gonbayes"
)

func TestP(t *testing.T) {
	tests := []struct {
		name       string
		categories []string
		dataset    map[string]string
		document   string
		want       map[string]float64
	}{
		{
			name:       "example",
			categories: []string{"経済", "IT", "エンタメ"},
			dataset: map[string]string{
				"Price":            "経済",
				"Insurance":        "経済",
				"Python and Price": "IT",
				"Python and Go":    "IT",
				"Marvel":           "エンタメ",
			},
			document: "Python and Go",
			want:     map[string]float64{"IT": 0.05, "エンタメ": 0, "経済": 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := gonbayes.NewClassifier(tt.categories, 0)
			for s, v := range tt.dataset {
				c.Train(v, s)
			}
			got := c.P(tt.document)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want = %v, got = %v\n", tt.want, got)
			}
		})
	}
}
