package x3dh_core

type X3DHInitBundle struct {
	// Identity Key
	IK X3DHPublicIK
	// Signed Pre Key
	SPK X3DHPublicSPK
	// One Time Pre Keys
	OtpSet []X3DHPublicOTP
}

type X3DHExpandeBundle struct {
	OtpSet []X3DHPublicOTP
}
