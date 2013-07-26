package main

import (
	"testing"
)

func TestGenerateGraph(t *testing.T) {
	buf := []byte("https://chart.googleapis.com/chart?cht=gv&chs=800x800&chl=digraph{")

	for _, e := range elements {
		for _, l := range e.Links {
			buf = append(append(append(append(buf, ' '), e.Name...), "->"...), l.Name...)
		}
	}

	t.Log(string(append(buf, " }"...)))
}

func TestElements(t *testing.T) {
	reverse := make(map[*element]Element)
	for i, e := range elements {
		reverse[e] = Element(i)
	}

	for i := Element(0); i < elementCount; i++ {
		if elements[i] == nil {
			t.Errorf("elements[%d] is nil", i)
		} else {
			var name string
			if len(elements[i].Name) == 0 {
				t.Errorf("elements[%d] has no name", i)
			} else {
				name = " (" + elements[i].Name + ")"
			}
			if len(elements[i].Links) == 0 {
				t.Errorf("elements[%d]%s has no links", i, name)
			}
			for _, l := range elements[i].Links {
				if l == elements[i] {
					t.Errorf("elements[%d]%s has a link to itself", i, name)
				}
				if _, ok := reverse[l]; !ok {
					t.Errorf("elements[%d]%s has an invalid link", i, name)
				}
			}
			for j := i + 1; j < elementCount; j++ {
				if elements[i] == elements[j] {
					t.Errorf("elements[%d]%s is a duplicate of elements[%d]", i, name, j)
				} else if elements[j] != nil {
					if elements[i].Name == elements[j].Name && name != "" {
						t.Errorf("elements[%d]%s is named the same thing as elements[%d]", i, name, j)
					}
				}
			}
		}
	}
}
