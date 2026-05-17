// sample.go — exercises common Go syntax features.
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const PI = 3.14159

// Circle is a closed plane curve with a fixed radius.
type Circle struct {
	Radius float64
}

// Area returns the area of the circle.
func (c Circle) Area() float64 {
	return PI * c.Radius * c.Radius
}

// Shape is implemented by anything with an Area.
type Shape interface {
	Area() float64
}

var errEmpty = errors.New("input is empty")

func greet(who string) string {
	if who == "" {
		who = "world"
	}
	return fmt.Sprintf("hello, %s", who)
}

func sumAreas(shapes []Shape) (float64, error) {
	if len(shapes) == 0 {
		return 0, errEmpty
	}
	var total float64
	for _, s := range shapes {
		total += s.Area()
	}
	return total, nil
}

func fetchUser(ctx context.Context, id int) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(50 * time.Millisecond):
		return map[string]any{"id": id, "name": "miru"}, nil
	}
}

func main() {
	fmt.Println(greet("miru"))

	shapes := []Shape{Circle{Radius: 1}, Circle{Radius: 5}}
	if total, err := sumAreas(shapes); err == nil {
		fmt.Printf("total area: %.2f\n", total)
	}

	if u, err := fetchUser(context.Background(), 42); err == nil {
		fmt.Println(strings.ToUpper(fmt.Sprint(u)))
	}
}
