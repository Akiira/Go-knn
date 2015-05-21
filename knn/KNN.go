// KNN
package knn

//TODO
//1. Change data set to k-d tree
//2. Remove hard coded values
//3. Implement cross folding
//4. Maybe have weighted voting function passed in as argument

import (
	"fmt"
	"sync"
)

//import (
////"KNN/knn/data"
////"KNN/knn/dist"
//)

type DataSet interface {
	Size() int
	GetEntry(int) interface{}
	GetClass(int) int
	Weight(index int) float64
}

type Neighbor struct {
	index    int
	distance float64
}

type Neighbors struct {
	neighbors []Neighbor
}

func NewNeighbors() *Neighbors {
	ns := new(Neighbors)
	ns.neighbors = make([]Neighbor, 0, 7) //TODO remove hardcoding
	return ns
}

func (n *Neighbors) ReplaceFurthest(dist float64, index int) {
	f_index, f_dist := n.Furthest()

	if f_dist > dist {
		n.neighbors[f_index].index = index
		n.neighbors[f_index].distance = dist
		//fmt.Println("\tReplacing furthest")
	} else {
		//fmt.Println("\tDidn't add")
	}
}

func (n *Neighbors) Furthest() (index int, dist float64) {
	var highest int = 0
	for index, elem := range n.neighbors {
		if elem.distance > n.neighbors[highest].distance {
			highest = index
		}
	}

	return highest, n.neighbors[highest].distance
}

func (n *Neighbors) AddIfCloser(dist float64, index int) {
	if len(n.neighbors) < K {
		n.neighbors = append(n.neighbors, Neighbor{index, dist})
	} else {
		n.ReplaceFurthest(dist, index)
	}
}

var (
	K           int
	data        DataSet
	folds       int
	GetDistance DistanceFunction
)

type DistanceFunction func(r1, r2 interface{}) float64

//TODO add argument for number of threads
func RunKNN(k int, dataSet DataSet, folds int, disFun DistanceFunction, numOfClasses int) {
	GetDistance = disFun
	data = dataSet
	K = k
	fmt.Println("Size of data set: ", data.Size())

	correct := 0
	total := 15000
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i <= total; i++ {
		//nghbrs := GetNeighbors(i) //TODO remove fixed index
		//fmt.Println("Neighbors: ", nghbrs)
		//class := Vote(nghbrs)
		//fmt.Println("Predicted: ", class, " for ", data.GetClass(i))

		go func(i2 int) {
			class := Vote(GetNeighbors(i2))
			if class == data.GetClass(i2) {
				correct++
			}
			wg.Done()
		}(i)

		if i%3 == 2 {
			wg.Wait()
			wg.Add(3)
		}
	}

	fmt.Println("Accuracy: ", float64(correct)/float64(total), ", Correct: ", correct, ", Total: ", total)
}

func Vote(neighbors Neighbors) (cuisine int) {
	var votes [8]float64 //TODO remove hard coding

	for _, neighbor := range neighbors.neighbors {
		dist := neighbor.distance
		//weight := data.Weight(neighbor.index)
		weight := 1.0
		//_ = weight * (1.0 / (dist * dist))
		votes[data.GetClass(neighbor.index)] += weight * (1.0 / (dist * dist))
		//votes[data.GetClass(neighbor.index)] += 1
	}

	var highest float64 = 0

	for index, vote := range votes {
		if vote > highest {
			cuisine = index
			highest = vote
		}
	}
	//fmt.Print("\tVotes: ", votes)
	return cuisine
}

func GetNeighbors(target int) Neighbors {
	neighbors := NewNeighbors()

	for other := 0; other < data.Size(); other++ {
		if other != target {
			dist := GetDistance(data.GetEntry(target), data.GetEntry(other))
			neighbors.AddIfCloser(dist, other)
		}
	}

	return *neighbors
}
