package main

import (
	"image"
	"image/color"
	"image/gif"
	"os"
)

const (
	pointSize = 12
	canvSize  = 400
)

var palette = []color.Color{color.White, color.Black, color.RGBA{0, 0, 255, 255}}

type GraphLine struct {
	Num   int
	Color uint8
}

func (g *GraphLine) Visit() {
	g.Color = 2
}

func (g *GraphLine) Visited() bool {
	return g.Color == 2
}

type Graph map[int][]GraphLine

func drawPoint(x, y int, img *image.Paletted, ci uint8) {
	for i := x - pointSize/2; i < x+pointSize/2; i++ {
		for j := y - pointSize/2; j < y+pointSize/2; j++ {
			img.SetColorIndex(i, j, ci)
		}
	}
}

func drawLine(x1, y1, x2, y2 int, img *image.Paletted, ci uint8) {
	var k float64 = float64(y2-y1) / float64(x2-x1)
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	for i := 0; i <= x2-x1; i++ {
		img.SetColorIndex(x1+i, y1+int(k*float64(i)), ci)
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	for i := 0; i <= y2-y1; i++ {
		img.SetColorIndex(x1+int(float64(i)/k), y1+i, ci)
	}
}

func (g *Graph) append(src, dst int) {
	val, ok := (*g)[src]
	if !ok {
		val = make([]GraphLine, 1)
		val[0] = GraphLine{dst, 1}
	} else {
		val = append(val, GraphLine{dst, 1})
	}
	(*g)[src] = val
}

func (g *Graph) drawVertext(v int, x1, x2 int, level int, img *image.Paletted, parentX, parentY int, color uint8) {
	rootX := x1 + (x2-x1-pointSize)/2
	rootY := 10 + 50*level

	childs := (*g)[v]

	drawPoint(rootX, rootY, img, color)
	if parentX != -1 {
		drawLine(parentX, parentY, rootX, rootY, img, color)
	}

	if len(childs) == 0 {
		return
	}
	stepX := (x2 - x1) / len(childs) // bug
	for i, vertex := range childs {
		childX1 := x1 + i*stepX
		childX2 := x1 + (i+1)*stepX
		g.drawVertext(vertex.Num, childX1, childX2, level+1, img, rootX, rootY, vertex.Color)
	}
}

func (g *Graph) draw(anim *gif.GIF) {

	rect := image.Rect(0, 0, canvSize, canvSize)
	img := image.NewPaletted(rect, palette)

	g.drawVertext(0, 0, canvSize, 0, img, -1, -1, 2)

	anim.Delay = append(anim.Delay, 50)
	anim.Image = append(anim.Image, img)
}

func dfs(g *Graph, anim *gif.GIF, v int) {
	g.draw(anim)

	childs := (*g)[v]

	for i, vertex := range childs {
		(*g)[v][i].Visit()
		dfs(g, anim, vertex.Num)
	}
}

func bfs(g *Graph, anim *gif.GIF, v int) {
	g.draw(anim)

	queue := make([]int, 1)
	queue[0] = v

	for {
		done := true
		for _, vertex := range queue {
			for i, child := range (*g)[vertex] {
				if (*g)[vertex][i].Visited() {
					continue
				}
				done = false

				(*g)[vertex][i].Visit()
				g.draw(anim)
				queue = append(queue, child.Num)
			}
		}
		if done {
			break
		}
	}

	childs := (*g)[v]

	for i, _ := range childs {
		(*g)[v][i].Visit()
		g.draw(anim)
	}

	for _, vertex := range childs {
		for i, _ := range (*g)[vertex.Num] {
			(*g)[vertex.Num][i].Visit()
			g.draw(anim)
		}
	}
}

func makeGraph() Graph {
	graph := make(Graph)
	graph.append(0, 1)
	graph.append(0, 2)

	graph.append(1, 3)
	graph.append(1, 4)

	graph.append(2, 20)
	graph.append(20, 200)
	graph.append(20, 201)

	graph.append(200, 2000)
	graph.append(200, 2001)

	graph.append(2001, 20010)

	graph.append(201, 2010)
	graph.append(201, 2011)

	graph.append(2010, 20100)
	graph.append(20100, 201000)
	graph.append(20100, 201001)
	graph.append(201001, 2010010)

	graph.append(4, 5)
	graph.append(4, 6)
	graph.append(4, 7)

	graph.append(5, 50)
	graph.append(5, 51)

	graph.append(51, 510)

	graph.append(510, 5100)

	return graph
}

func main() {
	{
		anim := gif.GIF{}

		graph := makeGraph()
		dfs(&graph, &anim, 0)

		f, _ := os.Create("dfs.gif")
		defer f.Close()
		gif.EncodeAll(f, &anim)
	}
	{
		anim := gif.GIF{}

		graph := makeGraph()
		bfs(&graph, &anim, 0)

		f, _ := os.Create("bfs.gif")
		defer f.Close()
		gif.EncodeAll(f, &anim)
	}
}
