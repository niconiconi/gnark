// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gnark DO NOT EDIT

package plonk_test

import (
	"testing"

	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/plonk"
	mockcommitment "github.com/consensys/gnark/crypto/polynomial/bw761/mock_commitment"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/internal/backend/bw761/cs"
	plonkbw761 "github.com/consensys/gnark/internal/backend/bw761/plonk"
	bw761witness "github.com/consensys/gnark/internal/backend/bw761/witness"
	"github.com/consensys/gnark/internal/backend/circuits"
	curve "github.com/consensys/gurvy/bw761"
)

func TestCircuits(t *testing.T) {
	for name, circuit := range circuits.Circuits {
		t.Run(name, func(t *testing.T) {
			assert := plonk.NewAssert(t)
			pcs, err := frontend.Compile(curve.ID, backend.PLONK, circuit.Circuit)
			assert.NoError(err)
			assert.SolvingSucceeded(pcs, circuit.Good)
			assert.SolvingFailed(pcs, circuit.Bad)
		})
	}
}

// TODO WIP -> once everything is clean move this to backend/plonk in assert
func TestProver(t *testing.T) {
	t.Skip("skip for bw761")

	for name, circuit := range circuits.Circuits {
		// name := "range"
		// circuit := circuits.Circuits[name]

		t.Run(name, func(t *testing.T) {

			assert := plonk.NewAssert(t)
			pcs, err := frontend.Compile(curve.ID, backend.PLONK, circuit.Circuit)
			assert.NoError(err)

			spr := pcs.(*cs.SparseR1CS)

			scheme := mockcommitment.Scheme{}
			wPublic := bw761witness.Witness{}
			wPublic.FromPublicAssignment(circuit.Good)
			publicData := plonkbw761.Setup(spr, &scheme, wPublic)

			// correct proof
			{
				wFull := bw761witness.Witness{}
				wFull.FromFullAssignment(circuit.Good)
				proof := plonkbw761.Prove(spr, publicData, wFull)

				v := plonkbw761.VerifyRaw(proof, publicData, wPublic)

				if !v {
					t.Fatal("Correct proof verification failed")
				}
			}

			//wrong proof
			{
				wFull := bw761witness.Witness{}
				wFull.FromFullAssignment(circuit.Bad)
				proof := plonkbw761.Prove(spr, publicData, wFull)

				v := plonkbw761.VerifyRaw(proof, publicData, wPublic)

				if v {
					t.Fatal("Wrong proof verification should have failed")
				}
			}
		})

	}
}
