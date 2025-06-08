package packages

import (
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

type AssocRespUserData struct {
	Length   base.ASNLength
	UserData structures.MDSEUserInfoStd
}

type AssociationResponse struct {
	SessionHeader
	// AssocRespSessionData
	// AssocRespPresentationHeader
	AssocRespUserData
	// AssocRespPresentationTrailer
}
