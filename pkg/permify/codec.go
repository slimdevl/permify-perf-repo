package permify

import (
	"encoding/json"
	"strings"
)

// Permify does not allow IDs to have values used in basic regex
// so org.1224 must transition to org_1234, then back out.
// the simplest way to do this is json custom marshalling.
// this way it is transparent to our application and just works.

const DecodedSeparator = "."
const EncodedSeparator = "_"

// encodeID replaces DecodedSeparator with EncodedSeparator
func encodeID(input string) string {
	return strings.ReplaceAll(input, DecodedSeparator, EncodedSeparator)
}

// decodeID replaces EncodedSeparator with DecodedSeparator
func decodeID(input string) string {
	return strings.ReplaceAll(input, EncodedSeparator, DecodedSeparator)
}

// encodeIDs replaces DecodedSeparator with EncodedSeparator in a slice of strings
func encodeIDs(ids []string) []string {
	encodedIDs := make([]string, len(ids))
	for i, id := range ids {
		encodedIDs[i] = strings.ReplaceAll(id, DecodedSeparator, EncodedSeparator)
	}
	return encodedIDs
}

// decodeIDs replaces EncodedSeparator with DecodedSeparator in a slice of strings
func decodeIDs(ids []string) []string {
	decodedIDs := make([]string, len(ids))
	for i, id := range ids {
		decodedIDs[i] = strings.ReplaceAll(id, EncodedSeparator, DecodedSeparator)
	}
	return decodedIDs
}

func (e *Entity) MarshalJSON() ([]byte, error) {
	e.Id = encodeID(e.Id)
	type Alias Entity
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	type Alias Entity
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	e.Id = decodeID(e.Id)
	return nil
}

func (s *Subject) MarshalJSON() ([]byte, error) {
	s.Id = encodeID(s.Id)
	type Alias Subject
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *Subject) UnmarshalJSON(data []byte) error {
	type Alias Subject
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Id = decodeID(s.Id)
	return nil
}

func (l *leaf) MarshalJSON() ([]byte, error) {
	l.Id = encodeID(l.Id)
	type Alias leaf
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

func (l *leaf) UnmarshalJSON(data []byte) error {
	type Alias leaf
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	l.Id = decodeID(l.Id)
	return nil
}

func (e *EntityIDSet) MarshalJSON() ([]byte, error) {
	e.Ids = encodeIDs(e.Ids)
	type Alias EntityIDSet
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

func (e *EntityIDSet) UnmarshalJSON(data []byte) error {
	type Alias EntityIDSet
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	e.Ids = decodeIDs(e.Ids)
	return nil
}

func (e *LookupRelationshipRequest) MarshalJSON() ([]byte, error) {
	e.Subject.Id = encodeID(e.Subject.Id)
	type Alias LookupRelationshipRequest
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

func (e *LookupRelationshipRequest) UnmarshalJSON(data []byte) error {
	type Alias LookupRelationshipRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	e.Subject.Id = decodeID(e.Subject.Id)
	return nil
}

func (e *LookupRelationshipResponse) MarshalJSON() ([]byte, error) {
	e.EntityIDs = encodeIDs(e.EntityIDs)
	type Alias LookupRelationshipResponse
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

func (e *LookupRelationshipResponse) UnmarshalJSON(data []byte) error {
	type Alias LookupRelationshipResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	e.EntityIDs = decodeIDs(e.EntityIDs)
	return nil
}

func (s *SubjectIDSet) MarshalJSON() ([]byte, error) {
	s.Ids = encodeIDs(s.Ids)
	type Alias SubjectIDSet
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *SubjectIDSet) UnmarshalJSON(data []byte) error {
	type Alias SubjectIDSet
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Ids = decodeIDs(s.Ids)
	return nil
}
