package types

// Lookup describes subscription lookup statuses
type Lookup string

// LookupStruct contains the Lookup value
type LookupStruct struct {
	Value *Lookup
}

const (
	// FoundLookup found
	FoundLookup Lookup = "FOUND"
	// NotFoundLookup not found
	NotFoundLookup Lookup = "NOT_FOUND"
	// ExpiredLookup expired
	ExpiredLookup Lookup = "EXPIRED"
)

// NewFoundLookup returns FoundLookup
func NewFoundLookup() LookupStruct {
	v := FoundLookup
	return LookupStruct{
		Value: &v,
	}
}

// NewExpiredLookup returns ExpiredLookup
func NewExpiredLookup() LookupStruct {
	v := ExpiredLookup
	return LookupStruct{
		Value: &v,
	}
}

// NewNotFoundLookup returns NotFoundLookup
func NewNotFoundLookup() LookupStruct {
	v := NotFoundLookup
	return LookupStruct{
		Value: &v,
	}
}
