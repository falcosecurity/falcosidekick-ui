package models

type Outputs []string

var outputs *Outputs

func CreateOutputs() *Outputs {
	outputs = new(Outputs)
	return outputs
}

func GetOutputs() *Outputs {
	return outputs
}

func (outputs *Outputs) Update(o []string) {
	*outputs = o
}
