package anomaly

// Package anomaly implements basic anomaly scenarios that can be triggered by
// profiles of gnbsim. The scenarios intentionally craft abnormal NAS or
// signalling messages in order to exercise error handling paths on the 5G core
// network.

import (
	crand "crypto/rand"
	mrand "math/rand"
	"time"

	"github.com/omec-project/gnbsim/common"
	simuectx "github.com/omec-project/gnbsim/simue/context"
	nasMessage "github.com/omec-project/nas/nasMessage"
	"github.com/omec-project/nas/nasTestpacket"
	"github.com/omec-project/nas/nasType"
)

// Trigger executes the anomaly scenario identified by name.
// Unknown names default to the "malformed-nas" scenario.
func Trigger(ue *simuectx.SimUe, name string) {
	switch name {
	case "malformed-nas":
		malformedNAS(ue)
	case "service-before-registration":
		serviceBeforeRegistration(ue)
	case "duplicate-registration":
		duplicateRegistration(ue)
	case "auth-response-before-request":
		authResponseBeforeRequest(ue)
	case "security-complete-before-command":
		securityCompleteBeforeCommand(ue)
	case "deregistration-before-registration":
		deregistrationBeforeRegistration(ue)
	case "empty-nas":
		emptyNAS(ue)
	default:
		malformedNAS(ue)
	}
}

// malformedNAS sends a NAS message with random payload. The payload does not
// follow NAS encoding rules and is expected to be rejected by the core.
func malformedNAS(ue *simuectx.SimUe) {
	l := 20
	b := make([]byte, l)
	if _, err := crand.Read(b); err != nil {
		// rand.Read from crypto may fail; fallback to math/rand
		for i := range b {
			b[i] = byte(mrand.Intn(256))
		}
	}
	msg := &common.UuMessage{}
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	msg.NasPdus = common.NasPduList{b}
	// Send directly to gNB without any NAS protections
	ue.WriteGnbUeChan <- msg
}

// serviceBeforeRegistration sends a Service Request before any Registration.
// The real UE never performed initial registration which should trigger error
// handling on the core network.
func serviceBeforeRegistration(ue *simuectx.SimUe) {
	msg := &common.UeMessage{}
	msg.Event = common.SERVICE_REQUEST_EVENT
	// Directly send to Real UE which forwards to the network
	ue.WriteRealUeChan <- msg
	// Allow some time for message to be processed before returning
	time.Sleep(10 * time.Millisecond)
}

// duplicateRegistration spams the core with two Registration Requests back to back.
func duplicateRegistration(ue *simuectx.SimUe) {
	msg := &common.UeMessage{}
	msg.Event = common.REG_REQUEST_EVENT
	ue.WriteRealUeChan <- msg
	time.Sleep(10 * time.Millisecond)
	ue.WriteRealUeChan <- msg
	time.Sleep(10 * time.Millisecond)
}

// authResponseBeforeRequest sends an Authentication Response without any prior challenge.
func authResponseBeforeRequest(ue *simuectx.SimUe) {
	msg := &common.UeMessage{}
	msg.Event = common.AUTH_RESPONSE_EVENT
	ue.WriteRealUeChan <- msg
	time.Sleep(10 * time.Millisecond)
}

// securityCompleteBeforeCommand emits a Security Mode Complete before the command.
func securityCompleteBeforeCommand(ue *simuectx.SimUe) {
	msg := &common.UeMessage{}
	msg.Event = common.SEC_MOD_COMPLETE_EVENT
	ue.WriteRealUeChan <- msg
	time.Sleep(10 * time.Millisecond)
}

// deregistrationBeforeRegistration issues a Deregistration Request prior to registration.
func deregistrationBeforeRegistration(ue *simuectx.SimUe) {
	mobileIdentity := nasType.MobileIdentity5GS{}
	nasPdu := nasTestpacket.GetDeregistrationRequest(nasMessage.AccessType3GPP, 0, 0, mobileIdentity)
	msg := &common.UuMessage{}
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	msg.NasPdus = common.NasPduList{nasPdu}
	ue.WriteGnbUeChan <- msg
}

// emptyNAS sends a zero-length NAS payload to the core network.
func emptyNAS(ue *simuectx.SimUe) {
	msg := &common.UuMessage{}
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	msg.NasPdus = common.NasPduList{[]byte{}}
	ue.WriteGnbUeChan <- msg
}
