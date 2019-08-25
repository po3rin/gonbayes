package gonbayes_test

import (
	"reflect"
	"testing"

	"github.com/po3rin/gonbayes"
)

func TestIsStopWords(t *testing.T) {
	tests := []struct {
		name string
		word string
		want bool
	}{
		{
			name: "stop word",
			word: "you",
			want: true,
		},
		{
			name: "stop word",
			word: "great",
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := gonbayes.IsStopWords(tt.word)
			if got != tt.want {
				t.Errorf("want = %v, got = %v\n", tt.want, got)
			}
		})
	}
}

func TestStem(t *testing.T) {
	tests := []struct {
		name string
		word string
		want string
	}{
		{
			name: "stop word",
			word: "swims",
			want: "swim",
		},
		{
			name: "stop word",
			word: "ceiling",
			want: "ceil",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := gonbayes.Stem(tt.word)
			if got != tt.want {
				t.Errorf("want = %v, got = %v\n", tt.want, got)
			}
		})
	}
}

func TestClean(t *testing.T) {
	tests := []struct {
		name string
		word string
		want string
	}{
		{
			name: "clean",
			word: "You are Best friend!!!!!(^^）",
			want: "you are best friend",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := gonbayes.Clean(tt.word)
			if got != tt.want {
				t.Errorf("want = %v, got = %v\n", tt.want, got)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name string
		word string
		want map[string]int
	}{
		{
			name: "count words",
			word: "I say hello. You say hi",
			want: map[string]int{
				"say":   2,
				"hello": 1,
				"hi":    1,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := gonbayes.CountWords(tt.word)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want = %v, got = %v\n", tt.want, got)
			}
		})
	}
}
