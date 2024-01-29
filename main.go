package main

import (
	r "groupietracker/routeur"
	t "groupietracker/temps"
)

func main() {
	t.InitTemplate()
	r.InitServe()
}
