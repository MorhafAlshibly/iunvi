package conversion

import (
	"errors"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
)

func StringToUniqueIdentifier(value string) (mssql.UniqueIdentifier, error) {
	guid, err := uuid.Parse(value)
	if err != nil {
		return mssql.UniqueIdentifier{}, err
	}
	guidBytes, err := guid.MarshalBinary()
	if err != nil {
		return mssql.UniqueIdentifier{}, err
	}
	if len(guidBytes) != 16 {
		return mssql.UniqueIdentifier{}, errors.New("invalid GUID length")
	}
	return mssql.UniqueIdentifier(guidBytes), nil
}
