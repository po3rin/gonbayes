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
	}{}

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
