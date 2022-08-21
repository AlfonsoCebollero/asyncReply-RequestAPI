package interfaces

type IWorkflowManager interface {
	StartWorkflow(string, string, int) (string, string, error)
	RetrieveWFStatus(string) (int, bool, error)
}
