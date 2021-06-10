package action

type action interface {
	Run() error
}
