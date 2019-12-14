package main

import "fmt"

type Path struct {
	pattern string
}

func (p Path) String() string {
	return p.pattern
}

// Line adds a relative line to the path
func (p Path) Line(dx, dy float32) Path {
	return Path{p.pattern + fmt.Sprintf(" l %f %f", dx, dy)}
}

// LineTo adds an absolute line to the path
func (p Path) LineTo(x, y float32) Path {
	return Path{p.pattern + fmt.Sprintf(" L %f %f", x, y)}
}

// Hor adds a relative horizontal line to the path
func (p Path) Hor(dx float32) Path {
	return Path{p.pattern + fmt.Sprintf(" h %f", dx)}
}

// HorTo adds an absolute horizontal line to the path
func (p Path) HorTo(x float32) Path {
	return Path{p.pattern + fmt.Sprintf(" H %f", x)}
}

// Ver adds a relative vertical line to the path
func (p Path) Ver(dy float32) Path {
	return Path{p.pattern + fmt.Sprintf(" v %f", dy)}
}

// VerTo adds an absolute vertical line to the path
func (p Path) VerTo(y float32) Path {
	return Path{p.pattern + fmt.Sprintf(" V %f", y)}
}

// Move moves the cursor relatively
func (p Path) Move(dx, dy float32) Path {
	return Path{p.pattern + fmt.Sprintf(" m %f %f", dx, dy)}
}

// MoveTo moves the cursor absolutely
func (p Path) MoveTo(x, y float32) Path {
	return Path{p.pattern + fmt.Sprintf(" M %f %f", x, y)}
}

// Arc adds a relative arc to the path
func (p Path) Arc(rx, ry, dx, dy, rot float32, large bool, sweep bool) Path {
	l := 0
	if large {
		l = 1
	}
	s := 0
	if sweep {
		s = 1
	}
	return Path{p.pattern + fmt.Sprintf(" a %f %f %f %d %d %f %f", rx, ry, rot, l, s, dx, dy)}
}

// ArcTo adds an absolute arc to the path
func (p Path) ArcTo(rx, ry, x, y, rot float32, large bool, sweep bool) Path {
	l := 0
	if large {
		l = 1
	}
	s := 0
	if sweep {
		s = 1
	}
	return Path{p.pattern + fmt.Sprintf(" A %f %f %f %d %d %f %f", rx, ry, rot, l, s, x, y)}
}
