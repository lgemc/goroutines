package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func plotMessages(messages []pipeMessage) {
	// Create some example data points
	data := make(plotter.XYs, len(messages))
	for i, message := range messages {
		data[i].X = float64(message.xOriginal)
		data[i].Y = float64(message.y)
	}

	// Create a new plot
	p := plot.New()

	// Set plot title and axis labels
	p.Title.Text = "XY Plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a scatter plotter and set its style
	s, err := plotter.NewScatter(data)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = plotutil.Color(0)

	// Add the scatter plotter to the plot
	p.Add(s)

	// Save the plot to a file
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "xy_plot.png"); err != nil {
		panic(err)
	}
}
