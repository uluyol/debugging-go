package main

func setTrue(x *bool) {
	*x = true
}

func main() {
	x := false
	go func() {
		setTrue(&x)
	}()
	for !x {
		// spin
	}
}
