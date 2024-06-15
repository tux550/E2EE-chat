package x3dh

import (
	"go.step.sm/crypto/x25519"
)

type X3DHKeyBundle struct {
	// Identity Key
	ik X3DHFullIK
	// Signed Pre Key
	spk X3DHFullSPK
	// One Time Pre Keys
	otp_set []X3DHFullOTP
	// One Time Pre Key ID Counter
	otp_counter int
}

type X3DHPublicKeyBundle struct {
	// Identity Key
	ik X3DHPublicIK
	// Signed Pre Key
	spk X3DHPublicSPK
	// One Time Pre Keys
	otp_set []X3DHPublicOTP
}

func GenereateKeyBundle() (*X3DHKeyBundle, error) {
	// Generate identity key
	ik, err := GenerateFullIK()
	if err != nil {
		return nil, err
	}
	// Generate signed pre key
	spk, err := GenerateFullSPK(ik.IdentityKey)
	if err != nil {
		return nil, err
	}
	// Initalize
	kb := &X3DHKeyBundle{
		ik:          *ik,
		spk:         *spk,
		otp_set:     make([]X3DHFullOTP, 0),
		otp_counter: 0,
	}
	// Generate one time pre keys
	for i := 0; i < 5; i++ {
		err := kb.GenerateOneTimePreKey()
		if err != nil {
			return nil, err
		}
	}
	// Return
	return kb, nil
}

func (kb *X3DHKeyBundle) GenerateOneTimePreKey() error {
	// Generate one time pre key
	otp, err := GenerateFullOTP(kb.otp_counter)
	if err != nil {
		return err
	}
	// Increment the counter
	kb.otp_counter++
	// Append the one time pre key
	kb.otp_set = append(kb.otp_set, *otp)
	return nil
}

func (kb *X3DHKeyBundle) PublicBundle() *X3DHPublicKeyBundle {
	// Create a new public key bundle
	pkb := &X3DHPublicKeyBundle{
		ik:  *kb.ik.PublicIK(),
		spk: *kb.spk.PublicSPK(),
	}
	// Add one time pre keys
	for _, otp := range kb.otp_set {
		pkb.otp_set = append(pkb.otp_set, *otp.PublicOTP())
	}
	return pkb
}

func (kb *X3DHPublicKeyBundle) Validate() bool {
	// Validate the signed pre key
	valid := x25519.Verify(kb.ik.IdentityKey, kb.spk.SignedPreKey, kb.spk.SignedPreKeySignature)
	return valid
}
