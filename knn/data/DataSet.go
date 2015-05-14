// DataSet
package data

type DataSet interface {
	Size() int
	GetEntry(int) interface{}
	GetClass(int) int
	//Randomize()
	Weight(index int) float64
}
