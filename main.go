package main

import (
	"fmt"
	"os"
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
	allPaths           [][]string
	startRoom, endRoom string
	nombreOfAnts       int
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

// func AfficherGraphe(G *Graphe) {
// 	fmt.Println("Liste d'adjacence du sommet :")
// 	for sommet, adj := range G.listeAdjac {
// 		temp := adj.head
// 		fmt.Printf("%s ->", sommet)
// 		for temp != nil {
// 			fmt.Printf(" %s", temp.valeur)
// 			temp = temp.suivant
// 		}
// 		fmt.Println()
// 	}
// }

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

func distributeAnts(paths [][]string, nombreOfAnts int) [][]int {
	fmt.Println(len(paths), nombreOfAnts)
	antDistribution := make([][]int, len(paths))
	for i := 0; i < nombreOfAnts; i++ {
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

func ParseInput(fileName string) (*Graphe, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
	}
	fmt.Println(string(file))
	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	edges := make([]string, 0)
	var nbAnts string
	for i, line := range lines {
		switch {
		case i == 0:
			nbAnts = line
		case line == "##start":
			startRoom = strings.Split(lines[i+1], " ")[0]
		case line == "##end":
			endRoom = strings.Split(lines[i+1], " ")[0]
		case strings.Contains(line, "-"):
			edges = append(edges, line)
		}
	}

	nombreOfAnts, _ = strconv.Atoi(nbAnts)
	G := InitGraphe(nombreOfAnts)
	for _, edge := range edges {
		nodes := strings.Split(edge, "-")
		if len(nodes) == 2 {
			AjouterArete(G, nodes[0], nodes[1])
		}
	}
	return G, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <input_file>")
		return
	}
	G, err := ParseInput(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}

	path := []string{}
	G.findAllPath(startRoom, path)
	// fmt.Println("All Paths Found:")
	// for _, path := range allPaths {
	// 	fmt.Println(path)
	// }

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})
	// fmt.Println(allPaths)
	filterNonOverlappingPaths()
	allPaths = newSet(allPaths)
	// fmt.Println("Meilleurs chemins trouvés :")
	// for _, path := range allPaths {
	// 	fmt.Println(path)
	// }
	// AfficherGraphe(G)
	antDistribution := distributeAnts(allPaths, nombreOfAnts)
	simulateAntMovement(allPaths, antDistribution)
}
