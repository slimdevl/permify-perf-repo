package permify

import (
	"encoding/json"
	"testing"
)

const (
	expectedOrgEncoding = `{"type":"type","id":"org_1224"}`
)

func TestEncodeID(t *testing.T) {
	input := "org.1224"
	expected := "org_1224"
	result := encodeID(input)

	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}

func TestDecodeID(t *testing.T) {
	input := "org_1224"
	expected := "org.1224"
	result := decodeID(input)

	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}

func TestEncodeIDs(t *testing.T) {
	input := []string{"org.1224", "abc.5678"}
	expected := []string{"org_1224", "abc_5678"}
	result := encodeIDs(input)

	if len(result) != len(expected) {
		t.Errorf("Expected length: %d, Got length: %d", len(expected), len(result))
	}

	for i, r := range result {
		if r != expected[i] {
			t.Errorf("Expected: %s, Got: %s", expected[i], r)
		}
	}
}

func TestDecodeIDs(t *testing.T) {
	input := []string{"org_1224", "abc_5678"}
	expected := []string{"org.1224", "abc.5678"}
	result := decodeIDs(input)

	if len(result) != len(expected) {
		t.Errorf("Expected length: %d, Got length: %d", len(expected), len(result))
	}

	for i, r := range result {
		if r != expected[i] {
			t.Errorf("Expected: %s, Got: %s", expected[i], r)
		}
	}
}

func TestEntityMarshalJSON(t *testing.T) {
	entity := &Entity{Type: "type", Id: "org.1224"}
	expected := expectedOrgEncoding
	result, err := json.Marshal(entity)

	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected: %s, Got: %s", expected, string(result))
	}
}

func TestEntityUnmarshalJSON(t *testing.T) {
	jsonStr := expectedOrgEncoding
	expected := &Entity{Type: "type", Id: "org.1224"}
	var entity Entity
	err := json.Unmarshal([]byte(jsonStr), &entity)

	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	if entity.Type != expected.Type || entity.Id != expected.Id {
		t.Errorf("Expected: %+v, Got: %+v", expected, entity)
	}
}

func TestSubjectMarshalJSON(t *testing.T) {
	subject := &Subject{Type: "type", Id: "org.1224"}
	expected := expectedOrgEncoding
	result, err := json.Marshal(subject)

	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected: %s, Got: %s", expected, string(result))
	}
}

func TestSubjectUnmarshalJSON(t *testing.T) {
	jsonStr := expectedOrgEncoding
	expected := &Subject{Type: "type", Id: "org.1224"}
	var subject Subject
	err := json.Unmarshal([]byte(jsonStr), &subject)

	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	if subject.Type != expected.Type || subject.Id != expected.Id {
		t.Errorf("Expected: %+v, Got: %+v", expected, subject)
	}
}

func TestLeafMarshalJSON(t *testing.T) {
	leaf := &leaf{Type: "type", Id: "org.1224", Relation: "rel"}
	expected := `{"type":"type","id":"org_1224","relation":"rel"}`
	result, err := json.Marshal(leaf)

	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected: %s, Got: %s", expected, string(result))
	}
}

func TestLeafUnmarshalJSON(t *testing.T) {
	jsonStr := `{"type":"type","id":"org_1224","relation":"rel"}`
	expected := &leaf{Type: "type", Id: "org.1224", Relation: "rel"}
	var l leaf
	err := json.Unmarshal([]byte(jsonStr), &l)

	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	if l.Type != expected.Type || l.Id != expected.Id || l.Relation != expected.Relation {
		t.Errorf("Expected: %+v, Got: %+v", expected, l)
	}
}

func TestEntityIDSetMarshalJSON(t *testing.T) {
	entityIDSet := &EntityIDSet{Type: "type", Ids: []string{"org.1224", "abc.5678"}}
	expected := `{"type":"type","ids":["org_1224","abc_5678"]}`
	result, err := json.Marshal(entityIDSet)

	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected: %s, Got: %s", expected, string(result))
	}
}

func TestEntityIDSetUnmarshalJSON(t *testing.T) {
	jsonStr := `{"type":"type","ids":["org_1224","abc_5678"]}`
	expected := &EntityIDSet{Type: "type", Ids: []string{"org.1224", "abc.5678"}}
	var entityIDSet EntityIDSet
	err := json.Unmarshal([]byte(jsonStr), &entityIDSet)

	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	if entityIDSet.Type != expected.Type || len(entityIDSet.Ids) != len(expected.Ids) {
		t.Errorf("Expected: %+v, Got: %+v", expected, entityIDSet)
	}
	for i, id := range entityIDSet.Ids {
		if id != expected.Ids[i] {
			t.Errorf("Expected: %s, Got: %s", expected.Ids[i], id)
		}
	}
}

func TestLookupRelationshipResponseMarshalJSON(t *testing.T) {
	lookupResponse := &LookupRelationshipResponse{
		ErrorResponse: &ErrorResponse{Code: 404, Message: "Not Found", Details: []string{"Detail1", "Detail2"}},
		EntityIDs:     []string{"org.1224", "abc.5678"},
	}
	expected := `{"code":404,"message":"Not Found","details":["Detail1","Detail2"],"entity_ids":["org_1224","abc_5678"]}`
	result, err := json.Marshal(lookupResponse)

	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected: %s, Got: %s", expected, string(result))
	}
}

func TestLookupRelationshipResponseUnmarshalJSON(t *testing.T) {
	jsonStr := `{"code":404,"message":"Not Found","details":["Detail1","Detail2"],"entity_ids":["org_1224","abc_5678"]}`
	expected := &LookupRelationshipResponse{
		ErrorResponse: &ErrorResponse{Code: 404, Message: "Not Found", Details: []string{"Detail1", "Detail2"}},
		EntityIDs:     []string{"org.1224", "abc.5678"},
	}
	var lookupResponse LookupRelationshipResponse
	err := json.Unmarshal([]byte(jsonStr), &lookupResponse)

	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	if lookupResponse.Code != expected.Code || lookupResponse.Message != expected.Message || len(lookupResponse.Details) != len(expected.Details) {
		t.Errorf("Expected: %+v, Got: %+v", expected, lookupResponse)
	}
	for i, detail := range lookupResponse.Details {
		if detail != expected.Details[i] {
			t.Errorf("Expected: %s, Got: %s", expected.Details[i], detail)
		}
	}
	if len(lookupResponse.EntityIDs) != len(expected.EntityIDs) {
		t.Errorf("Expected entity IDs length: %d, Got: %d", len(expected.EntityIDs), len(lookupResponse.EntityIDs))
	}
	for i, id := range lookupResponse.EntityIDs {
		if id != expected.EntityIDs[i] {
			t.Errorf("Expected: %s, Got: %s", expected.EntityIDs[i], id)
		}
	}
}

func TestSubjectIDSetMarshalJSON(t *testing.T) {
	subjectIDSet := &SubjectIDSet{Type: "type", Ids: []string{"org.1224", "abc.5678"}, Relation: "rel"}
	expected := `{"type":"type","ids":["org_1224","abc_5678"],"relation":"rel"}`
	result, err := json.Marshal(subjectIDSet)

	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected: %s, Got: %s", expected, string(result))
	}
}

func TestSubjectIDSetUnmarshalJSON(t *testing.T) {
	jsonStr := `{"type":"type","ids":["org_1224","abc_5678"],"relation":"rel"}`
	expected := &SubjectIDSet{Type: "type", Ids: []string{"org.1224", "abc.5678"}, Relation: "rel"}
	var subjectIDSet SubjectIDSet
	err := json.Unmarshal([]byte(jsonStr), &subjectIDSet)

	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	if subjectIDSet.Type != expected.Type || len(subjectIDSet.Ids) != len(expected.Ids) || subjectIDSet.Relation != expected.Relation {
		t.Errorf("Expected: %+v, Got: %+v", expected, subjectIDSet)
	}
	for i, id := range subjectIDSet.Ids {
		if id != expected.Ids[i] {
			t.Errorf("Expected: %s, Got: %s", expected.Ids[i], id)
		}
	}
}
