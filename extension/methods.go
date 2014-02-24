package extension

// For loops through each leaf node of a TestGroup.
// path contains the path from root to leaf.
func (g *TestGroup) For(visit func(path TestGroups)) {
	g.each(&groupStack{}, visit)
}

func (g *TestGroup) each(s *groupStack, visit func(path TestGroups)) {
	s.push(g)
	defer s.pop()
	for _, child := range g.Children {
		child.each(s, visit)
	}
	if len(g.Children) == 0 {
		visit(s.a)
	}
}
