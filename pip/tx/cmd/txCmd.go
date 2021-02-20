package cmd

type TxCmd interface {
	Commit() error
}
