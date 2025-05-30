package packages

import (
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

type SinglePollDataResult struct {
	SPpdu
	ROapdus
	RORSapdu
	ActionResult
	PollMdibDataReply
}
