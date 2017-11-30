package consumer

type Consumer interface {
	SendReindexMessage(string) error
	SendExportMessage(string) error
}
