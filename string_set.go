package main

type StringSet struct {
    items map[string]bool
}

func(s *StringSet) Add(t string) *StringSet {
    if s.items == nil {
        s.items = make(map[string]bool)
    }

    _, ok := s.items[t]
    if !ok {
        s.items[t] = true
    }

    return s
} 

func(s *StringSet) Clear() {
    s.items = make(map[string]bool)
}

func (s *StringSet) Delete(item string) bool {
    _, ok := s.items[item]
    if ok {
        delete(s.items, item)
    }

    return ok
}

func (s *StringSet) Has(item string) bool {
    _, ok := s.items[item]

    return ok
}

func (s *StringSet) Strings() []string {
    items := []string{}
    for i := range s.items {
        items = append(items, i)
    }

    return items
}

func (s *StringSet) Size() int {
    return len(s.items)
} 

