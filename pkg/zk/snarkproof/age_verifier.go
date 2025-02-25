package snarkproof

import (
	"fmt"
	"os"
	"time"

	"github.com/alejoacosta74/go-logger"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

func VerifyAgeProof(proofFileStr string, vkFileStr string) error {

	// load proof from file
	var proof groth16.Proof

	file, err := os.Open(proofFileStr)
	if err != nil {
		return fmt.Errorf("Failed to open proof file: %v", err)
	}
	logger.Debugf("Proof file %s opened successfully\n", proofFileStr)
	defer file.Close()

	// Create a concrete proof type that implements the interface
	proofData := groth16.NewProof(ecc.BN254)
	_, err = proofData.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("Failed to read proof from file: %v", err)
	}
	logger.Debugf("Proof read successfully from file %s\n", proofFileStr)
	proof = proofData // Assign the concrete type to the interface

	// load verification key from file
	vk := groth16.NewVerifyingKey(ecc.BN254)
	file, err = os.Open(vkFileStr)
	if err != nil {
		logger.Fatalf("Failed to open verification key file: %v", err)
	}
	logger.Debugf("Verification key file %s opened successfully\n", vkFileStr)
	defer file.Close()

	_, err = vk.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("Failed to read verification key from file: %v", err)
	}
	logger.Debugf("Verification key read successfully from file %s\n", vkFileStr)

	// get current year
	currentYear := time.Now().Year()

	// create witness for verification
	publicAssignment := AgeCircuit{CurrentYear: currentYear}
	publicWitness, err := frontend.NewWitness(&publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return fmt.Errorf("Failed to create public witness: %v", err)
	}

	// verify zk-SNARKproof
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		return fmt.Errorf("Failed to verify proof: %v", err)
	} else {
		logger.Debugf("✅ Proof verification result: %v\n", true)
	}
	return nil

}
