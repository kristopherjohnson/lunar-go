// Translation of
// <http://www.cs.brandeis.edu/~storer/LunarLander/LunarLander/LunarLanderListing.jpg>
// by Jim Storer from FOCAL to Go.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Global variables from original FOCAL code
//
// a - Altitude (miles)
// g - Gravity
// i - Intermediate altitude (miles)
// j - Intermediate velocity (miles/sec)
// k - Fuel rate (lbs/sec)
// l - Elapsed time (sec)
// m - Total weight (lbs)
// n - Empty weight (lbs, Note: m - n is remaining fuel weight)
// s - Time elapsed in current 10-second turn (sec)
// t - Time remaining in current 10-second turn (sec)
// v - Downward speed (miles/sec)
// w - Temporary working variable
// z - Thrust per pound of fuel burned
var (
	a, g, i, j, k, l, m, n, s, t, v, w, z float64
	echoInput                             bool
	inputReader                           *bufio.Reader
)

func main() {
	// Initialize the persistent input reader
	inputReader = bufio.NewReader(os.Stdin)

	if len(os.Args) > 1 {
		// If --echo is present, then write all input back to standard output.
		// (This is useful for testing with files as redirected input.)
		if os.Args[1] == "--echo" {
			echoInput = true
		}
	}

	fmt.Println("CONTROL CALLING LUNAR MODULE. MANUAL CONTROL IS NECESSARY")
	fmt.Println("YOU MAY RESET FUEL RATE K EACH 10 SECS TO 0 OR ANY VALUE")
	fmt.Println("BETWEEN 8 & 200 LBS/SEC. YOU'VE 16000 LBS FUEL. ESTIMATED")
	fmt.Print("FREE FALL IMPACT TIME-120 SECS. CAPSULE WEIGHT-32500 LBS\n\n\n")

	for {
		playGame()

		fmt.Println("\n\n\nTRY AGAIN?")
		if !acceptYesOrNo() {
			break
		}
	}

	fmt.Print("CONTROL OUT\n\n\n")
}

func playGame() {
	// 01.20 in original FOCAL code
	fmt.Print("FIRST RADAR CHECK COMING UP\n\n\n")
	fmt.Println("COMMENCE LANDING PROCEDURE")
	fmt.Println("TIME,SECS   ALTITUDE,MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE")

	a = 120
	v = 1
	m = 32500
	n = 16500
	g = .001
	z = 1.8
	l = 0

startTurn: // 02.10 in original FOCAL code
	fmt.Printf("%7.0f%16.0f%7.0f%15.2f%12.1f      ",
		l,
		math.Trunc(a),
		5280*(a-math.Trunc(a)),
		3600*v,
		m-n)

promptForK:
	fmt.Print("K=:")
	var err error
	k, err = acceptFloat()
	if err != nil || k < 0 || ((0 < k) && (k < 8)) || k > 200 {
		fmt.Print("NOT POSSIBLE")
		for x := 1; x <= 51; x++ {
			fmt.Print(".")
		}
		goto promptForK
	}

	t = 10

turnLoop:
	for { // 03.10 in original FOCAL code
		if m-n < .001 {
			goto fuelOut
		}

		if t < .001 {
			goto startTurn
		}

		s = t

		if n+s*k-m > 0 {
			s = (m - n) / k
		}

		applyThrust()

		if i <= 0 {
			goto loopUntilOnTheMoon
		}

		if (v > 0) && (j < 0) {
			for { // 08.10 in original FOCAL code
				// FOCAL-to-Go gotcha: In FOCAL, multiplication has a higher
				// precedence than division. In Go, they have the same
				// precedence and are evaluated left-to-right. So the
				// original FOCAL subexpression `m * g / z * k` can't be
				// copied as-is into Go: `z * k` has to be parenthesized to
				// get the same result.
				w = (1 - m*g/(z*k)) / 2
				s = m*v/(z*k*(w+math.Sqrt(w*w+v/z))) + 0.05
				applyThrust()
				if i <= 0 {
					goto loopUntilOnTheMoon
				}
				updateLanderState()
				if -j < 0 {
					goto turnLoop
				}
				if v <= 0 {
					goto turnLoop
				}
			}
		}

		updateLanderState()
	}

loopUntilOnTheMoon: // 07.10 in original FOCAL code
	for s >= .005 {
		s = 2 * a / (v + math.Sqrt(v*v+2*a*(g-z*k/m)))
		applyThrust()
		updateLanderState()
	}
	goto onTheMoon

fuelOut: // 04.10 in original FOCAL code
	fmt.Printf("FUEL OUT AT %8.2f SECS\n", l)
	s = (math.Sqrt(v*v+2*a*g) - v) / g
	v += g * s
	l += s

onTheMoon: // 05.10 in original FOCAL code
	fmt.Printf("ON THE MOON AT %8.2f SECS\n", l)
	w = 3600 * v
	fmt.Printf("IMPACT VELOCITY OF %8.2f M.P.H.\n", w)
	fmt.Printf("FUEL LEFT: %8.2f LBS\n", m-n)
	if w <= 1 {
		fmt.Println("PERFECT LANDING !-(LUCKY)")
	} else if w <= 10 {
		fmt.Println("GOOD LANDING-(COULD BE BETTER)")
	} else if w <= 22 {
		fmt.Println("CONGRATULATIONS ON A POOR LANDING")
	} else if w <= 40 {
		fmt.Println("CRAFT DAMAGE. GOOD LUCK")
	} else if w <= 60 {
		fmt.Println("CRASH LANDING-YOU'VE 5 HRS OXYGEN")
	} else {
		fmt.Println("SORRY,BUT THERE WERE NO SURVIVORS-YOU BLEW IT!")
		fmt.Printf("IN FACT YOU BLASTED A NEW LUNAR CRATER %8.2f FT. DEEP\n", w*.277777)
	}
}

// Subroutine at line 06.10 in original FOCAL code
func updateLanderState() {
	l += s
	t -= s
	m -= s * k
	a = i
	v = j
}

// Subroutine at line 09.10 in original FOCAL code
func applyThrust() {
	Q := s * k / m
	Q2 := Q * Q
	Q3 := Q2 * Q
	Q4 := Q3 * Q
	Q5 := Q4 * Q

	j = v + g*s + z*(-Q-Q2/2-Q3/3-Q4/4-Q5/5)
	i = a - g*s*s/2 - v*s + z*s*(Q/2+Q2/6+Q3/12+Q4/20+Q5/30)
}

// Read a floating-point value from stdin.
// Returns the parsed float64 value and nil error on success,
// or returns 0 and an error if input did not contain a valid number.
// Exits on EOF or other failure to read input.
func acceptFloat() (float64, error) {
	line := acceptLine()
	return strconv.ParseFloat(strings.TrimSpace(line), 64)
}

// Reads input and returns true if it starts with 'Y' or 'y', false otherwise.
// This matches the behavior of the original C/BASIC version which treats
// any non-Y input as "no" rather than reprompting.
// If unable to read input, exits.
func acceptYesOrNo() bool {
	fmt.Print("(ANS. YES OR NO):")
	line := acceptLine()

	line = strings.TrimSpace(line)
	if len(line) > 0 {
		switch line[0] {
		case 'y', 'Y':
			return true
		}
	}
	return false
}

// Reads a line of input.
// If unable to read input, exits the program instead of returning.
func acceptLine() string {
	line, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nEND OF INPUT")
		os.Exit(1)
	}

	// Remove the trailing newline character
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	// Also remove carriage return on Windows
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}

	if echoInput {
		fmt.Println(line)
	}

	return line
}
