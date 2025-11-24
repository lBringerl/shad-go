//go:build !solution

package hogwarts

func inverseGraph(prereqs map[string][]string) map[string][]string {
	graph := make(map[string][]string, 0)

	for k, reqs := range prereqs {
		for _, v := range reqs {
			graph[v] = append(graph[v], k)
		}
	}

	return graph
}

func getStartPoints(graph map[string][]string) []string {
	reqsCounter := make(map[string]int, 0)
	for course, reqs := range graph {
		reqsCounter[course] = len(reqs)
		for _, req := range reqs {
			if _, exists := reqsCounter[req]; !exists {
				reqsCounter[req] = 0
			}
		}
	}

	startPoints := make([]string, 0)

	for k, v := range reqsCounter {
		if v == 0 {
			startPoints = append(startPoints, k)
		}
	}

	return startPoints
}

func walkGraphNoDeps(graph map[string][]string, startPoint string) {
	visited := make(map[string]struct{}, 0)
	recursionStack := make(map[string]struct{}, 0)
	nextPoints := []string{startPoint}

	for len(nextPoints) != 0 {
		parent := nextPoints[len(nextPoints)-1]
		recursionStack[parent] = struct{}{}

		stackSize := len(nextPoints)
		for _, req := range graph[parent] {
			if _, exists := recursionStack[req]; exists {
				panic("circular dependence")
			}
			if _, exists := visited[req]; exists {
				continue
			}
			nextPoints = append(nextPoints, req)
			recursionStack[req] = struct{}{}
		}
		if stackSize == len(nextPoints) {
			nextPoints = nextPoints[:len(nextPoints)-1]
			delete(recursionStack, parent)
			visited[parent] = struct{}{}
		}
	}
}

func walkGraphBFS(graph map[string][]string, reqs map[string][]string, startPoints []string) []string {
	visited := make(map[string]struct{}, 0)
	nextPoints := startPoints
	courses := make([]string, 0)

	for len(nextPoints) != 0 {
		parent := nextPoints[0]
		visited[parent] = struct{}{}
		courses = append(courses, parent)
		nextPoints = nextPoints[1:]

		for _, nextPoint := range graph[parent] {
			if _, exists := visited[nextPoint]; exists {
				continue
			}

			allReqsVisisted := true
			for _, req := range reqs[nextPoint] {
				if _, exists := visited[req]; !exists {
					allReqsVisisted = false
					break
				}
			}

			if allReqsVisisted {
				nextPoints = append(nextPoints, nextPoint)
			}
		}
	}

	return courses
}

func GetCourseList(prereqs map[string][]string) []string {
	startPoints := getStartPoints(prereqs)
	if len(startPoints) == 0 {
		panic("no start points!")
	}

	graph := inverseGraph(prereqs)

	for _, startPoint := range startPoints {
		walkGraphNoDeps(graph, startPoint)
	}

	courses := walkGraphBFS(graph, prereqs, startPoints)

	return courses
}
