package errors

type ErrorType int

const (
	ErrSetRedisCache ErrorType = iota
	ErrKeyNotExists
	ErrRetrievingKey
	ErrUnknownError
)

func (e ErrorType) String() string {
	switch e {
	case ErrSetRedisCache:
		return "Error setting value in Redis Cache"
	case ErrKeyNotExists:
		return "Specified key %s does not exist in redis cache"
	case ErrUnknownError:
		return "Unknown error occuring while performing your request"
	case ErrRetrievingKey:
		return "Error retrieving key %s from redis cache"
	default:
		return ""
	}
}
