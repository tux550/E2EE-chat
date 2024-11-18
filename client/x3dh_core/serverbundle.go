package x3dh_core

type X3DHClientBundle struct {
	// Identity Key
	IK X3DHPublicIK `json:"identity_key"`
	// Signed Pre Key
	SPK X3DHPublicSPK `json:"signed_pre_key"`
	// One Time Pre Keys
	OtpSet []X3DHPublicOTP `json:"one_time_pre_keys"`
}

type X3DHExpandeBundle struct {
	OtpSet []X3DHPublicOTP
}
