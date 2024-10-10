package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	valeur  string
	suivant *Node
}

type listeAdjac struct {
	head *Node
}

type Graphe struct {
	nombreSommets int
	listeAdjac    map[string]*listeAdjac
}

var (
	allPaths                         [][]string
	nombreOfAnts, startRoom, endRoom string
)

func InitNode(valeur string) *Node {
	return &Node{
		valeur:  valeur,
		suivant: nil,
	}
}

func InitGraphe(V int) *Graphe {
	return &Graphe{
		nombreSommets: V,
		listeAdjac:    make(map[string]*listeAdjac),
	}
}

func AjouterArete(G *Graphe, src string, dest string) {
	newNode := InitNode(dest)
	if G.listeAdjac[src] == nil {
		G.listeAdjac[src] = &listeAdjac{}
	}
	newNode.suivant = G.listeAdjac[src].head
	G.listeAdjac[src].head = newNode

	newNode = InitNode(src)
	if G.listeAdjac[dest] == nil {
		G.listeAdjac[dest] = &listeAdjac{}
	}
	newNode.suivant = G.listeAdjac[dest].head
	G.listeAdjac[dest].head = newNode
}

func AfficherGraphe(G *Graphe) {
	fmt.Println("Liste d'adjacence du sommet :")
	for sommet, adj := range G.listeAdjac {
		temp := adj.head
		fmt.Printf("%s ->", sommet)
		for temp != nil {
			fmt.Printf(" %s", temp.valeur)
			temp = temp.suivant
		}
		fmt.Println()
	}
}

func (G *Graphe) findAllPath(startRoom string, Path []string) {
	Path = append(Path, startRoom)
	if startRoom == endRoom {
		allPaths = append(allPaths, append([]string(nil), Path...))
		return
	}
	temp := G.listeAdjac[startRoom].head
	for temp != nil {
		if !IsValidPath(temp.valeur, Path) {
			G.findAllPath(temp.valeur, Path)
		}
		temp = temp.suivant
	}
}

func IsValidPath(point string, visited []string) bool {
	for _, v := range visited {
		if v == point {
			return true
		}
	}
	return false
}

func filterNonOverlappingPaths() {
	var bestSolution [][]string
	for _, path := range allPaths {
		b := true
		currentSolution := [][]string{path}
		for _, otherPath := range allPaths {
			if !isOverlapping(currentSolution, otherPath) {
				currentSolution = append(currentSolution, otherPath)
				b = false
			}
		}
		if len(currentSolution) > len(bestSolution) || (len(currentSolution) == len(bestSolution) && b) {
			bestSolution = currentSolution
		}
		allPaths = bestSolution
	}
}

func isOverlapping(current [][]string, newPath []string) bool {
	for _, path := range current {
		for _, room := range path {
			if room != startRoom && room != endRoom {
				for _, newRoom := range newPath {
					if newRoom == room {
						return true
					}
				}
			}
		}
	}
	return false
}

func arrayToString(arr []string) string {
	return strings.Join(arr, " ")
}

func newSet(twoDArray [][]string) [][]string {
	set := make(map[string]bool)
	var result [][]string

	for _, arr := range twoDArray {
		key := arrayToString(arr)
		if !set[key] {
			set[key] = true
			result = append(result, arr)
		}
	}
	return result
}

func distributeAnts(paths [][]string, nbAnts int) [][]int {
	fmt.Println(len(paths))
	antDistribution := make([][]int, len(paths))
	for i := 0; i < nbAnts; i++ {
		antDistribution[i%len(paths)] = append(antDistribution[i%len(paths)], i+1)
	}
	return antDistribution
}

func simulateAntMovement(paths [][]string, antDistribution [][]int) {
	type AntPosition struct {
		ant  int
		path int
		step int
	}

	var antPositions []AntPosition
	for pathIndex, ants := range antDistribution {
		for _, ant := range ants {
			antPositions = append(antPositions, AntPosition{ant, pathIndex, 0})
		}
	}

	moveCount := 0
	for len(antPositions) > 0 {
		var moves []string
		var newPositions []AntPosition
		usedLinks := make(map[string]bool)

		for _, pos := range antPositions {
			if pos.step < len(paths[pos.path])-1 {
				currentRoom := paths[pos.path][pos.step]
				nextRoom := paths[pos.path][pos.step+1]
				link := currentRoom + "-" + nextRoom
				if !usedLinks[link] {
					moves = append(moves, fmt.Sprintf("L%d-%s", pos.ant, nextRoom))
					newPositions = append(newPositions, AntPosition{pos.ant, pos.path, pos.step + 1})
					usedLinks[link] = true
				} else {
					newPositions = append(newPositions, pos)
				}
			}
		}
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
		antPositions = newPositions
		moveCount++
	}
	fmt.Println("---------------------------------------------------------------------------------------")
	fmt.Printf("Nombre de mouvements : %d\n", moveCount-1)
}

func main() {
	fileName := `9
#rooms
##start
start 0 3
##end
end 10 1
C0 1 0
C1 2 0
C2 3 0
C3 4 0
I4 5 0
I5 6 0
A0 1 2
A1 2 1
A2 4 1
B0 1 4
B1 2 4
E2 6 4
D1 6 3
D2 7 3
D3 8 3
H4 4 2
H3 5 2
F2 6 2
F3 7 2
F4 8 2
G0 1 5
G1 2 5
G2 3 5
G3 4 5
G4 6 5
H3-F2
H3-H4
H4-A2
start-G0
G0-G1
G1-G2
G2-G3
G3-G4
G4-D3
start-A0
A0-A1
A0-D1
A1-A2
A1-B1
A2-end
A2-C3
start-B0
B0-B1
B1-E2
start-C0
C0-C1
C1-C2
C2-C3
C3-I4
D1-D2
D1-F2
D2-E2
D2-D3
D2-F3
D3-end
F2-F3
F3-F4
F4-end
I4-I5
I5-end
`

	lines := strings.Split(strings.TrimSpace(fileName), "\n")

	edges := make([]string, 0)

	for i, line := range lines {
		switch {
		case i == 0:
			nombreOfAnts = line
		case line == "##start":
			startRoom = strings.Split(lines[i+1], " ")[0]
		case line == "##end":
			endRoom = strings.Split(lines[i+1], " ")[0]
		case strings.Contains(line, "-"):
			edges = append(edges, line)
		}
	}

	nbAnts, _ := strconv.Atoi(nombreOfAnts)
	G := InitGraphe(nbAnts)
	for _, edge := range edges {
		nodes := strings.Split(edge, "-")
		if len(nodes) == 2 {
			AjouterArete(G, nodes[0], nodes[1])
		}
	}
	path := []string{}
	G.findAllPath(startRoom, path)
	fmt.Println("All Paths Found:")
	for _, path := range allPaths {
		fmt.Println(path)
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j]) 
	})
	fmt.Println(allPaths)
	filterNonOverlappingPaths()
	allPaths = newSet(allPaths)
	fmt.Println("Meilleurs chemins trouvés :")
	for _, path := range allPaths {
		fmt.Println(path)
	}
	AfficherGraphe(G)
	antDistribution := distributeAnts(allPaths, nbAnts)
	simulateAntMovement(allPaths, antDistribution)
}
