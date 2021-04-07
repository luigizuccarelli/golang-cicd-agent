package schema

import "time"

type Request struct {
	Id      string `json:"id,omitemptye"`
	Message string `json:"message"`
}

// Response schema
type Response struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResultsSchema struct {
	Project   string
	StartTime string
	StopTime  string
	Items     []ItemInfo
}

type ItemInfo struct {
	Name        string
	Fail        bool
	Pass        bool
	Time        int64
	DisplayTime string
	Log         string
}

type MapBinding struct {
	RepoUrl      string `json:"url"`
	RepoName     string `json:"name"`
	RepoHash     string `json:"hash"`
	ActorName    string `json:"actorname"`
	ActorEmail   string `json:"actoremail"`
	Message      string `json:"message"`
	TagVersion   string `json:"tagversion,omitempty"`
	InfraRepo    string `json:"infrarepo"`
	ActionDetail string `json:"actiondetail"`
	Action       string `json:"action"`
	Env          string `json:"env"`
	Workdir      string `json:"workdir"`
}

type GiteaSchema struct {
	Secret      string `json:"secret"`
	Action      string `json:"action"`
	Ref         string `json:"ref"`
	Before      string `json:"before"`
	After       string `json:"after"`
	CompareURL  string `json:"compare_url"`
	PullRequest struct {
		ID     int    `json:"id"`
		URL    string `json:"url"`
		Number int    `json:"number"`
		User   struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			FullName  string `json:"full_name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
			Language  string `json:"language"`
			IsAdmin   bool   `json:"is_admin"`
			Username  string `json:"username"`
		} `json:"user"`
		Title     string        `json:"title"`
		Body      string        `json:"body"`
		Labels    []interface{} `json:"labels"`
		Milestone interface{}   `json:"milestone"`
		Assignee  interface{}   `json:"assignee"`
		Assignees interface{}   `json:"assignees"`
		State     string        `json:"state"`
		//Comments       string        `json:"comments"`
		HTMLURL        string    `json:"html_url"`
		DiffURL        string    `json:"diff_url"`
		PatchURL       string    `json:"patch_url"`
		Mergeable      bool      `json:"mergeable"`
		Merged         bool      `json:"merged"`
		MergedAt       time.Time `json:"merged_at"`
		MergeCommitSha string    `json:"merge_commit_sha"`
		MergedBy       struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			FullName  string `json:"full_name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
			Language  string `json:"language"`
			IsAdmin   bool   `json:"is_admin"`
			Username  string `json:"username"`
		} `json:"merged_by"`
		Base struct {
			Label  string `json:"label"`
			Ref    string `json:"ref"`
			Sha    string `json:"sha"`
			RepoID int    `json:"repo_id"`
			Repo   struct {
				ID    int `json:"id"`
				Owner struct {
					ID        int    `json:"id"`
					Login     string `json:"login"`
					FullName  string `json:"full_name"`
					Email     string `json:"email"`
					AvatarURL string `json:"avatar_url"`
					Language  string `json:"language"`
					IsAdmin   bool   `json:"is_admin"`
					Username  string `json:"username"`
				} `json:"owner"`
				Name            string      `json:"name"`
				FullName        string      `json:"full_name"`
				Description     string      `json:"description"`
				Empty           bool        `json:"empty"`
				Private         bool        `json:"private"`
				Fork            bool        `json:"fork"`
				Parent          interface{} `json:"parent"`
				Mirror          bool        `json:"mirror"`
				Size            int         `json:"size"`
				HTMLURL         string      `json:"html_url"`
				SSHURL          string      `json:"ssh_url"`
				CloneURL        string      `json:"clone_url"`
				Website         string      `json:"website"`
				StarsCount      int         `json:"stars_count"`
				ForksCount      int         `json:"forks_count"`
				WatchersCount   int         `json:"watchers_count"`
				OpenIssuesCount int         `json:"open_issues_count"`
				DefaultBranch   string      `json:"default_branch"`
				Archived        bool        `json:"archived"`
				CreatedAt       time.Time   `json:"created_at"`
				UpdatedAt       time.Time   `json:"updated_at"`
				Permissions     struct {
					Admin bool `json:"admin"`
					Push  bool `json:"push"`
					Pull  bool `json:"pull"`
				} `json:"permissions"`
			} `json:"repo"`
		} `json:"base"`
		Head struct {
			Label  string `json:"label"`
			Ref    string `json:"ref"`
			Sha    string `json:"sha"`
			RepoID int    `json:"repo_id"`
			Repo   struct {
				ID    int `json:"id"`
				Owner struct {
					ID        int    `json:"id"`
					Login     string `json:"login"`
					FullName  string `json:"full_name"`
					Email     string `json:"email"`
					AvatarURL string `json:"avatar_url"`
					Language  string `json:"language"`
					IsAdmin   bool   `json:"is_admin"`
					Username  string `json:"username"`
				} `json:"owner"`
				Name            string      `json:"name"`
				FullName        string      `json:"full_name"`
				Description     string      `json:"description"`
				Empty           bool        `json:"empty"`
				Private         bool        `json:"private"`
				Fork            bool        `json:"fork"`
				Parent          interface{} `json:"parent"`
				Mirror          bool        `json:"mirror"`
				Size            int         `json:"size"`
				HTMLURL         string      `json:"html_url"`
				SSHURL          string      `json:"ssh_url"`
				CloneURL        string      `json:"clone_url"`
				Website         string      `json:"website"`
				StarsCount      int         `json:"stars_count"`
				ForksCount      int         `json:"forks_count"`
				WatchersCount   int         `json:"watchers_count"`
				OpenIssuesCount int         `json:"open_issues_count"`
				DefaultBranch   string      `json:"default_branch"`
				Archived        bool        `json:"archived"`
				CreatedAt       time.Time   `json:"created_at"`
				UpdatedAt       time.Time   `json:"updated_at"`
				Permissions     struct {
					Admin bool `json:"admin"`
					Push  bool `json:"push"`
					Pull  bool `json:"pull"`
				} `json:"permissions"`
			} `json:"repo"`
		} `json:"head"`
		MergeBase string      `json:"merge_base"`
		DueDate   interface{} `json:"due_date"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt time.Time   `json:"updated_at"`
		ClosedAt  interface{} `json:"closed_at"`
	} `json:"pull_request"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
		Author  struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Verification interface{} `json:"verification"`
		Timestamp    time.Time   `json:"timestamp"`
		Added        interface{} `json:"added"`
		Removed      interface{} `json:"removed"`
		Modified     interface{} `json:"modified"`
	} `json:"commits"`
	HeadCommit interface{} `json:"head_commit"`
	Release    struct {
		ID              int       `json:"id"`
		TagName         string    `json:"tag_name"`
		TargetCommitish string    `json:"target_commitish"`
		Name            string    `json:"name"`
		Body            string    `json:"body"`
		URL             string    `json:"url"`
		TarballURL      string    `json:"tarball_url"`
		ZipballURL      string    `json:"zipball_url"`
		Draft           bool      `json:"draft"`
		Prerelease      bool      `json:"prerelease"`
		CreatedAt       time.Time `json:"created_at"`
		PublishedAt     time.Time `json:"published_at"`
		Author          struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			FullName  string `json:"full_name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
			Language  string `json:"language"`
			IsAdmin   bool   `json:"is_admin"`
			Username  string `json:"username"`
		} `json:"author"`
		Assets []interface{} `json:"assets"`
	} `json:"release"`
	Repository struct {
		ID    int `json:"id"`
		Owner struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			FullName  string `json:"full_name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
			Language  string `json:"language"`
			IsAdmin   bool   `json:"is_admin"`
			Username  string `json:"username"`
		} `json:"owner"`
		Name            string      `json:"name"`
		FullName        string      `json:"full_name"`
		Description     string      `json:"description"`
		Empty           bool        `json:"empty"`
		Private         bool        `json:"private"`
		Fork            bool        `json:"fork"`
		Parent          interface{} `json:"parent"`
		Mirror          bool        `json:"mirror"`
		Size            int         `json:"size"`
		HTMLURL         string      `json:"html_url"`
		SSHURL          string      `json:"ssh_url"`
		CloneURL        string      `json:"clone_url"`
		Website         string      `json:"website"`
		StarsCount      int         `json:"stars_count"`
		ForksCount      int         `json:"forks_count"`
		WatchersCount   int         `json:"watchers_count"`
		OpenIssuesCount int         `json:"open_issues_count"`
		DefaultBranch   string      `json:"default_branch"`
		Archived        bool        `json:"archived"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		Permissions     struct {
			Admin bool `json:"admin"`
			Push  bool `json:"push"`
			Pull  bool `json:"pull"`
		} `json:"permissions"`
	} `json:"repository"`
	Pusher struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Language  string `json:"language"`
		IsAdmin   bool   `json:"is_admin"`
		Username  string `json:"username"`
	} `json:"pusher"`
	Sender struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Language  string `json:"language"`
		IsAdmin   bool   `json:"is_admin"`
		Username  string `json:"username"`
	} `json:"sender"`
}

type Pipeline struct {
	Id   int
	Task string
	Name string
	Log  bool
	Meta string
}
