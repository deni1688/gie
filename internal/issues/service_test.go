package issues

import (
	"errors"
	"reflect"
	"testing"
)

type mockGitProvider struct {
	repos []Repo
	issue Issue
	err   error
}

func (m *mockGitProvider) GetRepos() (*[]Repo, error) {
	return &m.repos, m.err
}

func (m *mockGitProvider) CreateIssue(repo *Repo, issue *Issue) error {
	issue.ID = m.issue.ID
	return m.err
}

type mockNotifier struct {
	err error
}

func (m mockNotifier) Notify(issues *[]Issue) error {
	return m.err
}

func TestNew(t *testing.T) {
	type args struct {
		gitProvider GitProvider
		notifier    Notifier
		prefix      string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "should return a new service",
			args: args{
				gitProvider: &mockGitProvider{},
				notifier:    &mockNotifier{},
				prefix:      "prefix",
			},
			want: &service{
				gitProvider: &mockGitProvider{},
				notifier:    &mockNotifier{},
				prefix:      "prefix",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.gitProvider, tt.args.notifier, tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceExtractIssues(t *testing.T) {
	type fields struct {
		gitProvider GitProvider
		notifier    Notifier
		prefix      string
	}
	type args struct {
		content *string
		source  *string
	}
	content := `// TODO: issue1
					func test() int { // TODO: issue2
						// some other content
                        return 1+1
					}`
	s := "/path/to/file"

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]Issue
		wantErr bool
	}{
		{
			name: "should return a list of issues",
			fields: fields{
				gitProvider: &mockGitProvider{},
				notifier:    &mockNotifier{},
				prefix:      "// TODO:",
			},
			args: args{
				content: &content,
				source:  &s,
			},
			want: &[]Issue{
				{
					ID:            0,
					Title:         "Issue1",
					Desc:          "Extracted from /path/to/file",
					Url:           "",
					ExtractedLine: "// TODO: issue1\n",
				},
				{
					ID:            0,
					Title:         "Issue2",
					Desc:          "Extracted from /path/to/file",
					Url:           "",
					ExtractedLine: "// TODO: issue2\n",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := service{
				gitProvider: tt.fields.gitProvider,
				notifier:    tt.fields.notifier,
				prefix:      tt.fields.prefix,
			}
			got, err := r.ExtractIssues(tt.args.content, tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractIssues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractIssues() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceFindRepoByName(t *testing.T) {
	type fields struct {
		gitProvider GitProvider
		notifier    Notifier
		prefix      string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Repo
		wantErr bool
	}{
		{
			name: "should return a repo if found",
			fields: fields{
				gitProvider: &mockGitProvider{
					repos: []Repo{{ID: 1, Name: "repo1", Owner: "owner1"}},
				},
				notifier: &mockNotifier{},
				prefix:   "prefix",
			},
			args: args{
				name: "repo1",
			},
			want:    &Repo{ID: 1, Name: "repo1", Owner: "owner1"},
			wantErr: false,
		},
		{
			name: "should return an error if not found",
			fields: fields{
				gitProvider: &mockGitProvider{
					repos: []Repo{{ID: 1, Name: "repo1", Owner: "owner1"}},
				},
				notifier: &mockNotifier{},
				prefix:   "prefix",
			},
			args: args{
				name: "repo2",
			},
			want:    &Repo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := service{
				gitProvider: tt.fields.gitProvider,
				notifier:    tt.fields.notifier,
				prefix:      tt.fields.prefix,
			}
			got, err := r.FindRepoByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRepoByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindRepoByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceGetUpdatedLine(t *testing.T) {
	type fields struct {
		gitProvider GitProvider
		notifier    Notifier
		prefix      string
	}
	type args struct {
		issue Issue
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "should return the extracted line updated with issue link",
			fields: fields{
				gitProvider: &mockGitProvider{},
				notifier:    &mockNotifier{},
				prefix:      "test",
			},
			args: args{
				issue: Issue{
					ID:            1212,
					Title:         "Make code better",
					Url:           "https://github.com/owner/repo/issues/1212",
					ExtractedLine: "// TODO: Make code better",
				},
			},
			want: "// TODO: Make code better -> closes https://github.com/owner/repo/issues/1212\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := service{
				gitProvider: tt.fields.gitProvider,
				notifier:    tt.fields.notifier,
				prefix:      tt.fields.prefix,
			}
			if got := r.GetUpdatedLine(tt.args.issue); got != tt.want {
				t.Errorf("GetUpdatedLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceNotify(t *testing.T) {
	type fields struct {
		gitProvider GitProvider
		notifier    Notifier
		prefix      string
	}
	type args struct {
		issues *[]Issue
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := service{
				gitProvider: tt.fields.gitProvider,
				notifier:    tt.fields.notifier,
				prefix:      tt.fields.prefix,
			}
			if err := r.Notify(tt.args.issues); (err != nil) != tt.wantErr {
				t.Errorf("Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceSubmitIssue(t *testing.T) {
	type fields struct {
		gitProvider GitProvider
		notifier    Notifier
		prefix      string
	}
	type args struct {
		repo  *Repo
		issue *Issue
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := service{
				gitProvider: tt.fields.gitProvider,
				notifier:    tt.fields.notifier,
				prefix:      tt.fields.prefix,
			}
			if err := r.SubmitIssue(tt.args.repo, tt.args.issue); (err != nil) != tt.wantErr {
				t.Errorf("SubmitIssue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceListRepos(t *testing.T) {
	mockRepos := []Repo{{ID: 1, Name: "repo1", Owner: "owner1"}}
	type fields struct {
		gitProvider GitProvider
		notifier    Notifier
		prefix      string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *[]Repo
		wantErr bool
	}{
		{
			name: "listRepos returns the repos from the git provider",
			fields: fields{
				gitProvider: &mockGitProvider{
					repos: mockRepos,
					issue: Issue{},
					err:   nil,
				},
				notifier: &mockNotifier{nil},
				prefix:   "// test",
			},
			want:    &mockRepos,
			wantErr: false,
		},
		{
			name: "listRepos returns an error if the git provider returns an error",
			fields: fields{
				gitProvider: &mockGitProvider{
					repos: []Repo{},
					issue: Issue{},
					err:   errors.New("mock error"),
				},
				notifier: &mockNotifier{nil},
				prefix:   "// test",
			},
			want:    &[]Repo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := service{
				gitProvider: tt.fields.gitProvider,
				notifier:    tt.fields.notifier,
				prefix:      tt.fields.prefix,
			}
			got, err := r.listRepos()
			if (err != nil) != tt.wantErr {
				t.Errorf("listRepos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listRepos() got = %v, want %v", got, tt.want)
			}
		})
	}
}
