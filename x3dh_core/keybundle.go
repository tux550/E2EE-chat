package x3dh_core

import (
	"go.step.sm/crypto/x25519"
)

type X3DHKeyBundle struct {
	// Identity Key
	IK X3DHPublicIK
	// Signed Pre Key
	SPK X3DHPublicSPK
	// One Time Pre Key
	OTP X3DHPublicOTP
}

func (kb *X3DHKeyBundle) Validate() bool {
	// Validate the signed pre key
	valid := x25519.Verify(kb.IK.IdentityKey, kb.SPK.SignedPreKey, kb.SPK.SignedPreKeySignature)
	return valid
}
