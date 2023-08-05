package main

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

//go:embed testdata/blockParser.json
var blockParserInput []byte

//go:embed testdata/sceneParser.json
var sceneParserInput []byte

func TestBlockParser(t *testing.T) {
	block := Block{}
	if err := json.Unmarshal(blockParserInput, &block); err != nil {
		t.Fatal(err)
	}

	expectedBlock := Block{pos: mgl32.Vec2{0, 80}, size: mgl32.Vec2{400, 20}}

	aPos, bPos := expectedBlock.pos, block.pos
	aSize, bSize := expectedBlock.size, block.size

	if aPos.X() != bPos.X() || aPos.Y() != bPos.Y() ||
		aSize.X() != bSize.X() || aSize.Y() != bSize.Y() {
		t.Fatalf("parsed block does not match")
	}
}

func TestSceneParser(t *testing.T) {
	scene := Scene{}
	if err := json.Unmarshal(sceneParserInput, &scene); err != nil {
		t.Fatal(err)
	}

	expectedBlocks := []Block{
		{pos: mgl32.Vec2{0, 80}, size: mgl32.Vec2{400, 20}},
		{pos: mgl32.Vec2{0, 0}, size: mgl32.Vec2{400, 20}},
		{pos: mgl32.Vec2{400, 20}, size: mgl32.Vec2{400, 20}},
		{pos: mgl32.Vec2{800, 40}, size: mgl32.Vec2{20, 100}},
		{pos: mgl32.Vec2{820, 170}, size: mgl32.Vec2{400, 20}},
	}

	if len(scene.Blocks) != len(expectedBlocks) {
		t.Fatalf("expected number of parsed block does not match, expected = %d, got = %d",
			len(expectedBlocks), len(scene.Blocks))
	}

	for i := 0; i < len(expectedBlocks); i += 1 {
		aPos, bPos := expectedBlocks[i].pos, scene.Blocks[i].pos
		aSize, bSize := expectedBlocks[i].size, scene.Blocks[i].size

		if aPos.X() != bPos.X() || aPos.Y() != bPos.Y() ||
			aSize.X() != bSize.X() || aSize.Y() != bSize.Y() {
			t.Fatal("parsed block value does not match")
		}
	}
}
