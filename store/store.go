package store

type Store interface {
	Write 
	Read
}

type Write interface {
	AddTransfer() error
}

type Read interface {

}