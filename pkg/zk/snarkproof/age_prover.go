package snarkproof

import (
	"fmt"
	"os"
	"time"

	"github.com/alejoacosta74/go-logger"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func GenerateAgeProof(userBirthYear int) error {
	currentYear := time.Now().Year()

	// Create witness: User proves their age ≥ threshold
	assignment := AgeCircuit{
		BirthYear:   userBirthYear, // User input
		CurrentYear: currentYear,   // Public input
	}

	// Compile the zk-SNARK circuit
	var circuit AgeCircuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return fmt.Errorf("Failed to compile R1CS: %v", err)
	}

	// Generate zk-SNARK proving and verification keys
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return fmt.Errorf("Failed to setup proving system: %v", err)
	}

	// Save the verification key for future use
	keyFile, err := os.Create(vkFileStr)
	if err != nil {
		return fmt.Errorf("Failed to create vk file: %v", err)
	}
	defer keyFile.Close()

	_, err = vk.WriteTo(keyFile) // Use WriteTo instead of gob encoding
	if err != nil {
		return fmt.Errorf("Failed to write vk to file: %v", err)
	}
	logger.Infof("✅ Verification key written successfully to file %s\n", vkFileStr)

	// Create private witness for proving
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return fmt.Errorf("Failed to create witness: %v", err)
	}

	// Generate zk-SNARK proof
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return fmt.Errorf("Failed to generate proof: %v", err)
	}

	// Save the proof to a file
	// saveToFile(proofFile, proof)
	proofFile, err := os.Create(proofFileStr)
	if err != nil {
		return fmt.Errorf("Failed to create proof file: %v", err)
	}
	defer proofFile.Close()

	_, err = proof.WriteTo(proofFile)
	if err != nil {
		return fmt.Errorf("Failed to write proof to file: %v", err)
	}

	logger.Debugf("✅ Proof successfully generated and saved in file %s \n", proofFileStr)

	return nil
}
