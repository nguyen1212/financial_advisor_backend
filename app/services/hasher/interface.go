package hasher

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

type I interface {
	Hash(string) []byte
}
