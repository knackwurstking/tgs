package opmanga

type PDF struct {
	name string
	data []byte
}

func NewPDF(name string, data []byte) *PDF {
	if data == nil {
		data = []byte{}
	}

	return &PDF{
		name: name,
		data: data,
	}
}

func (pdf *PDF) Name() string {
	return pdf.name
}

func (pdf *PDF) Data() []byte {
	return pdf.data
}
