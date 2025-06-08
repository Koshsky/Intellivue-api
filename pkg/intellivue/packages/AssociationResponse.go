package packages

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

type AssocRespUserData struct {
	Length   base.ASNLength
	UserData structures.MDSEUserInfoStd
}

type AssociationResponse struct {
	AssocRespSessionHeader       structures.SessionHeader
	AssocRespSessionData         structures.SessionData
	AssocRespPresentationHeader  structures.PresentationHeader
	AssocRespUserData            structures.UserData
	AssocRespPresentationTrailer []byte
}

func (a *AssociationResponse) Size() uint16 {
	return a.AssocRespSessionHeader.Size() +
		a.AssocRespSessionData.Size() +
		a.AssocRespPresentationHeader.Size() +
		a.AssocRespUserData.Size() +
		uint16(len(a.AssocRespPresentationTrailer))
}

func (a *AssociationResponse) UnmarshalBinary(r io.Reader) error {
	if a == nil {
		return fmt.Errorf("nil AssociationResponse receiver")
	}

	if err := a.AssocRespSessionHeader.UnmarshalBinary(r, base.AC_SPDU_SI); err != nil {
		return fmt.Errorf("failed to unmarshal AssocRespSessionHeader: %w", err)
	}

	if err := a.AssocRespSessionData.UnmarshalBinary(r, 14); err != nil {
		return fmt.Errorf("failed to unmarshal AssocRespSessionData: %w", err)
	}

	if err := a.AssocRespPresentationHeader.UnmarshalBinary(r, 101); err != nil {
		return fmt.Errorf("failed to unmarshal AssocRespPresentationHeader: %w", err)
	}

	if err := a.AssocRespUserData.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal AssocRespUserData: %w", err)
	}

	a.AssocRespPresentationTrailer = make([]byte, 16)
	length, err := io.ReadFull(r, a.AssocRespPresentationTrailer)
	if err != nil {
		return fmt.Errorf("failed to read AssocRespPresentationTrailer: %w", err)
	} else if length != 16 {
		return fmt.Errorf("failed to read AssocRespPresentationTrailer: read %d bytes, expected 16", length)
	}

	return nil
}

func (a *AssociationResponse) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<AssociationResponse>", indent)
	a.AssocRespSessionHeader.ShowInfo(indentationLevel + 1)
	a.AssocRespSessionData.ShowInfo(indentationLevel + 1)
	a.AssocRespPresentationHeader.ShowInfo(indentationLevel + 1)
	a.AssocRespUserData.ShowInfo(indentationLevel + 1)
	log.Printf("%s  AssocRespPresentationTrailer: 0x%X", indent, a.AssocRespPresentationTrailer)
}
