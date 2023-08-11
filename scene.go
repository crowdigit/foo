package main

type Scene struct {
	Blocks []Block `json:"blocks"`
}

func (s Scene) RenderBlocks(renderer RectRenderer) {
	for _, block := range s.Blocks {
		block.Render(renderer)
	}
}
