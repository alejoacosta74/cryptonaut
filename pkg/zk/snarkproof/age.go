package snarkproof

import (
	"github.com/consensys/gnark/frontend"
)

const (
	// default paths for generated proof and verification keys
	proofFileStr = "proof.data"
	vkFileStr    = "vk.data"

	// 🎯 Age threshold (set to 18 for this demo, can be changed)
	AgeThreshold = 18
)

// 🎯 zk-SNARK Circuit for Age Verification
type AgeCircuit struct {
	BirthYear   frontend.Variable `gnark:",private"` // Private input: User's birth year
	CurrentYear frontend.Variable `gnark:",public"`  // Public input: Known current year
}

// 🎯 Define the constraints: Ensures (CurrentYear - BirthYear) >= AgeThreshold
func (c *AgeCircuit) Define(api frontend.API) error {
	// Compute user's age: CurrentYear - BirthYear
	age := api.Sub(c.CurrentYear, c.BirthYear)

	// Assert that age >= AgeThreshold
	api.AssertIsLessOrEqual(frontend.Variable(AgeThreshold), age)

	return nil
}
