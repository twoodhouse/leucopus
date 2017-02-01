package truthTable

import (
	"math"
	"strconv"
	"strings"
)

var uidCounter = 0

type Link struct {
	Table          *TruthTable
	SourceLinks    []*Link //this is necessary for ordering correctly
	TargetLinks    []*Link
	Inputs         []int
	InputsSnapshot []int
	ExitLinkInputs []int
	ExitILink      bool
	Output         int
	OutputSnapshot int
	Uid            string
}

func NewLink(table *TruthTable, exitILink bool) *Link {
	inputs := make([]int, 0)
	for i := 0; i < table.Size; i++ {
		inputs = append(inputs, -1)
	}
	defaults := make([]int, 0)
	for i := 0; i < table.Size; i++ {
		defaults = append(defaults, 1)
	}
	uid := strconv.Itoa(uidCounter)
	uidCounter = uidCounter + 1
	var entity = Link{
		table,
		make([]*Link, table.Size),
		[]*Link{},
		inputs,
		inputs,
		inputs,
		exitILink,
		-1,
		-1,
		uid,
	}
	return &entity
}

func (l *Link) TakeSnapshot() {
	l.OutputSnapshot = l.Output
	for i, value := range l.Inputs {
		l.InputsSnapshot[i] = value
	}
}

func (l *Link) RestoreSnapshot() {
	l.Output = l.OutputSnapshot
	for i, value := range l.InputsSnapshot {
		l.Inputs[i] = value
	}
}

func (l *Link) FeedInputByLink(lnk *Link, val int) []*Link {
	var pos int
	for i, e := range l.SourceLinks {
		if e == lnk {
			pos = i
		}
	}
	l.Inputs[pos] = val
	inputsFull := true
	for _, e := range l.Inputs {
		if e == -1 {
			inputsFull = false
		}
	}
	if inputsFull && !l.ExitILink {
		l.Process()
		return l.Forward()
	}
	return []*Link{l}
}

func (l *Link) Process() {

	l.Output = l.Table.Process(l.Inputs)
	// print("processing ")
	// println(l.Uid)
	// println(l.Output)
	// l.Table.Print()
	// l.Print()
}

func (l *Link) Forward() []*Link {
	var unfinishedLinks []*Link

	for _, targetLink := range l.TargetLinks {
		//if target has all inputs already, clear them:
		needsCleared := true
		for _, input := range targetLink.Inputs {
			if input == -1 {
				needsCleared = false
			}
		}
		if needsCleared {
			targetLink.Inputs = make([]int, 0)
			for i := 0; i < targetLink.Table.Size; i++ {
				targetLink.Inputs = append(targetLink.Inputs, -1)
			}
		}

		intermediateUnfinishedLinks := targetLink.FeedInputByLink(l, l.Output) //this line actually does the work
		for _, newLink := range intermediateUnfinishedLinks {                  //now add all the unifinished links together
			existsAlready := false
			for _, refLink := range unfinishedLinks {
				if newLink == refLink {
					existsAlready = true
				}
			}
			if !existsAlready {
				unfinishedLinks = append(unfinishedLinks, newLink)
			}
		}
	}
	return unfinishedLinks
}

func AttachLinks(source *Link, target *Link, sourceNum int) {
	source.TargetLinks = append(source.TargetLinks, target)
	target.SourceLinks[sourceNum] = source
}

type TruthTable struct {
	Outputs           []int //2 is equivalent to a dash
	outputsStaging    []int
	Size              int
	LastProcessResult int
}

func New(outputs []int) *TruthTable {
	sze := size(outputs)
	var truthTable = TruthTable{
		outputs,
		outputs,
		sze,
		0,
	}
	return &truthTable
}

func NewEntryTable() *TruthTable {
	var truthTable = TruthTable{
		[]int{0, 1},
		[]int{0, 1},
		1,
		0,
	}
	return &truthTable
}

func (t *TruthTable) Process(inputs []int) int {
	/*
		Discussion: what happens when this function is fed a 2? Currently it seems that it outputs 0.
		If the function has 2s in its output and is fed 1s and 0s, then it will output a 2.
	*/
	if len(inputs) != t.Size {
		println("Warning: truth table process has wrong number of inputs for truth table")
	}
	for _, e := range inputs {
		if e == -1 {
			println("Error: You may have forgotten to set a starting value for an ExitILink")
		}
	}
	var inputsText []string
	for i := range inputs {
		number := inputs[i]
		text := strconv.Itoa(number)
		inputsText = append(inputsText, text)
	}
	location, _ := strconv.ParseInt(strings.Join(inputsText, ""), 2, 64)
	t.LastProcessResult = t.Outputs[location]
	return t.Outputs[location]
}

func (t *TruthTable) ReplaceInputValue(inputs []int, value int) {
	if len(inputs) != t.Size {
		println("Warning: truth table process has wrong number of inputs for truth table")
	}
	for _, e := range inputs {
		if e == -1 {
			println("Error: You may have forgotten to set a starting value for an ExitILink")
		}
	}
	var inputsText []string
	for i := range inputs {
		number := inputs[i]
		text := strconv.Itoa(number)
		inputsText = append(inputsText, text)
	}
	location, _ := strconv.ParseInt(strings.Join(inputsText, ""), 2, 64)
	t.Outputs[location] = value
}

func (t *TruthTable) Fill() {
	lastOutputs := make([]int, len(t.Outputs))
	same := false
	for same == false {
		for i, e := range t.Outputs {
			if e == 2 {
				t.outputsStaging[i] = t.fillSpot(i, []int{i})
			}
		}
		for i := range t.outputsStaging {
			t.Outputs[i] = t.outputsStaging[i]
		}
		same = true
		for i := range t.Outputs {
			if lastOutputs[i] != t.Outputs[i] {
				same = false
			}
		}
		lastOutputs = t.Outputs
	}
}

func (t *TruthTable) fillSpot(spot int, spotsComplete []int) int {
	//gather an experiment for each close relative
	n := int64(spot)
	binaryStrRow := strings.Split(strconv.FormatInt(n, 2), "")
	binaryIntRow := make([]int, t.Size)
	for i, e := range binaryStrRow {
		binaryIntRow[i], _ = strconv.Atoi(e)
	}
	var experiments [][]int
	for i := 0; i < t.Size; i++ {
		inputsCopy := make([]int, t.Size)
		copy(inputsCopy, binaryIntRow)
		if binaryIntRow[i] == 0 {
			inputsCopy[i] = 1
		} else if binaryIntRow[i] == 1 {
			inputsCopy[i] = 0
		}
		experiments = append(experiments, inputsCopy)
	}
	//now run experiments
	ones := 0
	zeroes := 0
	dashes := 0
	for _, experiment := range experiments {
		n := t.Process(experiment)
		if n == 0 {
			zeroes = zeroes + 1
		}
		if n == 1 {
			ones = ones + 1
		}
		if n == 2 {
			dashes = dashes + 1
		}
	}
	if zeroes > ones {
		return 0
	}
	if ones > zeroes {
		return 1
	}
	//in this case, it is a tie: Go deeper
	ones2 := 0
	zeroes2 := 0
	dashes2 := 0
	for _, experiment := range experiments {
		expText := ""
		for _, e := range experiment {
			expText = expText + strconv.Itoa(e)
		}
		location, _ := strconv.ParseInt(expText, 2, 64)
		oldLocation := false
		for _, e := range spotsComplete {
			if e == int(location) {
				oldLocation = true
			}
		}
		n2 := 2
		if oldLocation != true {
			n2 = t.fillSpot(int(location), append(spotsComplete, int(location)))
		}
		if n2 == 0 {
			zeroes2 = zeroes2 + 1
		}
		if n2 == 1 {
			ones2 = ones2 + 1
		}
		if n2 == 2 {
			dashes2 = dashes2 + 1
		}
	}
	if zeroes2 > ones2 {
		return 0
	}
	if ones2 > zeroes2 {
		return 1
	}
	return 2
}

func size(outputs []int) int {
	return int(math.Log2(float64(len(outputs))))
}

func (t *TruthTable) Print() {
	for _, e := range t.Outputs {
		print(e)
	}
	println()
}

func (l *Link) Print() {
	print("Inputs: ")
	for _, input := range l.Inputs {
		print(input)
		print(",")
	}
	println()
	print("Output: " + strconv.Itoa(l.Output))
	println()
}
