//go:build !solution

package retryupdate

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var (
		value      *string
		apiErr     *kvapi.APIError
		oldVersion = uuid.Nil
	)
GetLoop:
	for {
		resp, err := c.Get(&kvapi.GetRequest{
			Key: key,
		})
		if err != nil {
			if errors.Is(err, kvapi.ErrKeyNotFound) {
				value = nil
				oldVersion = uuid.Nil
			} else if errors.As(err, &apiErr) {
				switch apiErr.Err.(type) {
				case *kvapi.ConflictError:
					continue GetLoop
				case *kvapi.AuthError:
					return fmt.Errorf("c.Get: %w", err)
				default:
					continue GetLoop
				}
			}
		} else {
			value = &resp.Value
			oldVersion = resp.Version
		}

		generatedVersions := make(map[uuid.UUID]struct{})

		for {
			newValue, err := updateFn(value)
			if err != nil {
				return fmt.Errorf("updateFn: %w", err)
			}

			newVersion := uuid.Must(uuid.NewV4())
			generatedVersions[newVersion] = struct{}{}
			_, err = c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      newValue,
				OldVersion: oldVersion,
				NewVersion: newVersion,
			})
			if err != nil {
				if errors.Is(err, kvapi.ErrKeyNotFound) {
					value = nil
					oldVersion = uuid.Nil
					continue
				} else if errors.As(err, &apiErr) {
					switch apiErr.Err.(type) {
					case *kvapi.ConflictError:
						conflictErr := apiErr.Err.(*kvapi.ConflictError)
						_, exists := generatedVersions[conflictErr.ExpectedVersion]
						if exists {
							return nil
						} else {
							continue GetLoop
						}
					case *kvapi.AuthError:
						return fmt.Errorf("c.Get: %w", err)
					default:
						continue
					}
				}
			}

			return nil
		}
	}
}
