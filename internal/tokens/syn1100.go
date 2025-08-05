package tokens

import "fmt"

// HealthRecord represents encrypted healthcare data tied to an owner.
type HealthRecord struct {
	Owner  string
	Data   []byte
	Access map[string]bool
}

// SYN1100Token manages healthcare records with access control.
type SYN1100Token struct {
	records map[TokenID]*HealthRecord
}

// NewSYN1100Token creates an empty healthcare record store.
func NewSYN1100Token() *SYN1100Token {
	return &SYN1100Token{records: make(map[TokenID]*HealthRecord)}
}

// AddRecord stores a new healthcare record with the given ID.
func (t *SYN1100Token) AddRecord(id TokenID, owner string, data []byte) error {
	if _, exists := t.records[id]; exists {
		return fmt.Errorf("record exists")
	}
	t.records[id] = &HealthRecord{Owner: owner, Data: data, Access: map[string]bool{owner: true}}
	return nil
}

// GrantAccess allows a grantee to read the specified record.
func (t *SYN1100Token) GrantAccess(id TokenID, grantee string) error {
	rec, ok := t.records[id]
	if !ok {
		return fmt.Errorf("record not found")
	}
	rec.Access[grantee] = true
	return nil
}

// RevokeAccess revokes a previously granted permission.
func (t *SYN1100Token) RevokeAccess(id TokenID, grantee string) error {
	rec, ok := t.records[id]
	if !ok {
		return fmt.Errorf("record not found")
	}
	delete(rec.Access, grantee)
	return nil
}

// GetRecord returns the record data if the caller has access rights.
func (t *SYN1100Token) GetRecord(id TokenID, caller string) ([]byte, error) {
	rec, ok := t.records[id]
	if !ok {
		return nil, fmt.Errorf("record not found")
	}
	if !rec.Access[caller] {
		return nil, fmt.Errorf("access denied")
	}
	return rec.Data, nil
}
