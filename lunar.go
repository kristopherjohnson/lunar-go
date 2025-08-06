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

// Global variables
//
// A - Altitude (miles)
// G - Gravity
// I - Intermediate altitude (miles)
// J - Intermediate velocity (miles/sec)
// K - Fuel rate (lbs/sec)
// L - Elapsed time (sec)
// M - Total weight (lbs)
// N - Empty weight (lbs, Note: M - N is remaining fuel weight)
// S - Time elapsed in current 10-second turn (sec)
// T - Time remaining in current 10-second turn (sec)
// V - Downward speed (miles/sec)
// W - Temporary working variable
// Z - Thrust per pound of fuel burned
var (
	A, G, I, J, K, L, M, N, S, T, V, W, Z float64
	echoInput                             bool
)

func main() {
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
	fmt.Println("FREE FALL IMPACT TIME-120 SECS. CAPSULE WEIGHT-32500 LBS\n\n")

	for {
		playGame()

		fmt.Println("\n\n\nTRY AGAIN?")
		if !acceptYesOrNo() {
			break
		}
	}

	fmt.Println("CONTROL OUT\n\n")
}

func playGame() {
	// 01.20 in original FOCAL code
	fmt.Println("FIRST RADAR CHECK COMING UP\n\n")
	fmt.Println("COMMENCE LANDING PROCEDURE")
	fmt.Println("TIME,SECS   ALTITUDE,MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE")

	A = 120
	V = 1
	M = 32500
	N = 16500
	G = .001
	Z = 1.8
	L = 0

startTurn: // 02.10 in original FOCAL code
	fmt.Printf("%7.0f%16.0f%7.0f%15.2f%12.1f      ",
		L,
		math.Trunc(A),
		5280*(A-math.Trunc(A)),
		3600*V,
		M-N)

promptForK:
	fmt.Print("K=:")
	var err error
	K, err = acceptDouble()
	if err != nil || K < 0 || ((0 < K) && (K < 8)) || K > 200 {
		fmt.Print("NOT POSSIBLE")
		for x := 1; x <= 51; x++ {
			fmt.Print(".")
		}
		goto promptForK
	}

	T = 10

	for { // 03.10 in original FOCAL code
		if M-N < .001 {
			goto fuelOut
		}

		if T < .001 {
			goto startTurn
		}

		S = T

		if N+S*K-M > 0 {
			S = (M - N) / K
		}

		applyThrust()

		if I <= 0 {
			goto loopUntilOnTheMoon
		}

		if (V > 0) && (J < 0) {
			for { // 08.10 in original FOCAL code
				// FOCAL-to-Go gotcha: In FOCAL, multiplication has a higher
				// precedence than division. In Go, they have the same
				// precedence and are evaluated left-to-right. So the
				// original FOCAL subexpression `M * G / Z * K` can't be
				// copied as-is into Go: `Z * K` has to be parenthesized to
				// get the same result.
				W = (1 - M*G/(Z*K)) / 2
				S = M*V/(Z*K*(W+math.Sqrt(W*W+V/Z))) + 0.05
				applyThrust()
				if I <= 0 {
					goto loopUntilOnTheMoon
				}
				updateLanderState()
				if -J < 0 {
					goto startTurn
				}
				if V <= 0 {
					goto startTurn
				}
			}
		}

		updateLanderState()
	}

loopUntilOnTheMoon: // 07.10 in original FOCAL code
	for S >= .005 {
		S = 2 * A / (V + math.Sqrt(V*V+2*A*(G-Z*K/M)))
		applyThrust()
		updateLanderState()
	}
	goto onTheMoon

fuelOut: // 04.10 in original FOCAL code
	fmt.Printf("FUEL OUT AT %8.2f SECS\n", L)
	S = (math.Sqrt(V*V+2*A*G) - V) / G
	V += G * S
	L += S

onTheMoon: // 05.10 in original FOCAL code
	fmt.Printf("ON THE MOON AT %8.2f SECS\n", L)
	W = 3600 * V
	fmt.Printf("IMPACT VELOCITY OF %8.2f M.P.H.\n", W)
	fmt.Printf("FUEL LEFT: %8.2f LBS\n", M-N)
	if W <= 1 {
		fmt.Println("PERFECT LANDING !-(LUCKY)")
	} else if W <= 10 {
		fmt.Println("GOOD LANDING-(COULD BE BETTER)")
	} else if W <= 22 {
		fmt.Println("CONGRATULATIONS ON A POOR LANDING")
	} else if W <= 40 {
		fmt.Println("CRAFT DAMAGE. GOOD LUCK")
	} else if W <= 60 {
		fmt.Println("CRASH LANDING-YOU'VE 5 HRS OXYGEN")
	} else {
		fmt.Println("SORRY,BUT THERE WERE NO SURVIVORS-YOU BLEW IT!")
		fmt.Printf("IN FACT YOU BLASTED A NEW LUNAR CRATER %8.2f FT. DEEP\n", W*.277777)
	}
}

// Subroutine at line 06.10 in original FOCAL code
func updateLanderState() {
	L += S
	T -= S
	M -= S * K
	A = I
	V = J
}

// Subroutine at line 09.10 in original FOCAL code
func applyThrust() {
	Q := S * K / M
	Q2 := math.Pow(Q, 2)
	Q3 := math.Pow(Q, 3)
	Q4 := math.Pow(Q, 4)
	Q5 := math.Pow(Q, 5)

	J = V + G*S + Z*(-Q-Q2/2-Q3/3-Q4/4-Q5/5)
	I = A - G*S*S/2 - V*S + Z*S*(Q/2+Q2/6+Q3/12+Q4/20+Q5/30)
}

// Read a floating-point value from stdin.
// Returns the parsed float64 value and nil error on success,
// or returns 0 and an error if input did not contain a valid number.
// Exits on EOF or other failure to read input.
func acceptDouble() (float64, error) {
	line := acceptLine()

	f, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
	return f, err
}

// Reads input and returns true if it starts with 'Y' or 'y', or returns false if it
// starts with 'N' or 'n'.
// If input starts with none of those characters, prompt again.
// If unable to read input, exits.
func acceptYesOrNo() bool {
	for {
		fmt.Print("(ANS. YES OR NO):")
		line := acceptLine()

		line = strings.TrimSpace(line)
		if len(line) > 0 {
			switch line[0] {
			case 'y', 'Y':
				return true
			case 'n', 'N':
				return false
			}
		}
	}
}

// Reads a line of input.
// If unable to read input, exits.
func acceptLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "\nEND OF INPUT")
		os.Exit(1)
	}

	line := scanner.Text() + "\n"

	if echoInput {
		fmt.Print(line)
	}

	return strings.TrimSuffix(line, "\n")
}
