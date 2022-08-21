package tests

import (
	"fmt"
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/config"
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/useCases"
	"os"
	"path"
	"testing"
)

func init() {
	pwd, _ := os.Getwd()
	root := path.Join(path.Dir(pwd), "../../")
	os.Chdir(root)
	fmt.Println("Changed directory to: " + root)
	pwd, _ = os.Getwd()

	config.AppConfig.LoadConfiguration()
}

func Test_workflowCreationRequestValidation(t *testing.T) {
	type args struct {
		taskListName string
		duration     string
		name         string
		wfs          map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		want2   int
		wantErr bool
	}{
		{
			name: "Successful validation",
			args: args{
				taskListName: "",
				duration:     "35",
				name:         "waitingwf",
				wfs:          config.AppConfig.Cadence.Workflows,
			},
			want:    "asyncAPI",
			want1:   "WaitingWorkflow",
			want2:   35,
			wantErr: false,
		},
		{
			name: "Invalid duration",
			args: args{
				taskListName: "",
				duration:     "25",
				name:         "waitingwf",
				wfs:          config.AppConfig.Cadence.Workflows,
			},
			want:    "",
			want1:   "",
			want2:   0,
			wantErr: true,
		},
		{
			name: "Non existing workflow",
			args: args{
				taskListName: "",
				duration:     "35",
				name:         "nonExisting",
				wfs:          config.AppConfig.Cadence.Workflows,
			},
			want:    "",
			want1:   "nonExisting",
			want2:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := useCases.WorkflowCreationRequestValidation(tt.args.taskListName, tt.args.duration, tt.args.name, tt.args.wfs)
			if (err != nil) != tt.wantErr {
				t.Errorf("workflowCreationRequestValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("workflowCreationRequestValidation() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("workflowCreationRequestValidation() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("workflowCreationRequestValidation() got2 = %v, want %v", got2, tt.want2)
			}

		})
	}
}
