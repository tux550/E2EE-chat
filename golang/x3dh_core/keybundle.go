package x3dh_core

import (
	"fmt"

	"go.step.sm/crypto/x25519"
)

type X3DHKeyBundle struct {
	// Identity Key
	IK X3DHPublicIK `json:"identity_key"`
	// Signed Pre Key
	SPK X3DHPublicSPK `json:"signed_pre_key"`
	// One Time Pre Key
	OTP X3DHPublicOTP `json:"one_time_pre_key"`
}

func (kb *X3DHKeyBundle) Validate() bool {
	// Validate the signed pre key
	valid := x25519.Verify(kb.IK.IdentityKey, kb.SPK.SignedPreKey, kb.SPK.SignedPreKeySignature)
	return valid
}

func (kb *X3DHKeyBundle) DebugPrint() {
	fmt.Println("Identity Key:", kb.IK.IdentityKey)
	fmt.Println("Signed Pre Key:", kb.SPK.SignedPreKey)
	fmt.Println("OTP Key:", kb.OTP.OneTimePreKey)
}
