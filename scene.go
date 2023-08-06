package main

type Scene struct {
	Blocks []Block `json:"blocks"`
}

func (s Scene) RenderBlocks(renderer Renderer) {
	for _, block := range s.Blocks {
		block.Render(renderer)
	}
}
