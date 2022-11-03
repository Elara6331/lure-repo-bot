package types

import "time"

type PullRequestPayload struct {
	IsGitea bool   `json:"-"`
	Action  string `json:"action"`
	Number  int64  `json:"number"`
	Changes struct {
		Title struct {
			From string `json:"from"`
		} `json:"title"`
		Body struct {
			From string `json:"from"`
		} `json:"body"`
	} `json:"changes"`
	Label        Label        `json:"label"`
	PullRequest  PullRequest  `json:"pull_request"`
	Repository   Repository   `json:"repository"`
	Organization Organization `json:"organization"`
	Sender       User         `json:"sender"`
}

type Repository struct {
	ID               int64     `json:"id"`
	NodeID           string    `json:"node_id"`
	Name             string    `json:"name"`
	FullName         string    `json:"full_name"`
	Private          bool      `json:"private"`
	Owner            User      `json:"owner"`
	HTMLURL          string    `json:"html_url"`
	Description      string    `json:"description"`
	Fork             bool      `json:"fork"`
	URL              string    `json:"url"`
	ForksURL         string    `json:"forks_url"`
	KeysURL          string    `json:"keys_url"`
	CollaboratorsURL string    `json:"collaborators_url"`
	TeamsURL         string    `json:"teams_url"`
	HooksURL         string    `json:"hooks_url"`
	IssueEventsURL   string    `json:"issue_events_url"`
	EventsURL        string    `json:"events_url"`
	AssigneesURL     string    `json:"assignees_url"`
	BranchesURL      string    `json:"branches_url"`
	TagsURL          string    `json:"tags_url"`
	BlobsURL         string    `json:"blobs_url"`
	GitTagsURL       string    `json:"git_tags_url"`
	GitRefsURL       string    `json:"git_refs_url"`
	TreesURL         string    `json:"trees_url"`
	StatusesURL      string    `json:"statuses_url"`
	LanguagesURL     string    `json:"languages_url"`
	StargazersURL    string    `json:"stargazers_url"`
	ContributorsURL  string    `json:"contributors_url"`
	SubscribersURL   string    `json:"subscribers_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	CommitsURL       string    `json:"commits_url"`
	GitCommitsURL    string    `json:"git_commits_url"`
	CommentsURL      string    `json:"comments_url"`
	IssueCommentURL  string    `json:"issue_comment_url"`
	ContentsURL      string    `json:"contents_url"`
	CompareURL       string    `json:"compare_url"`
	MergesURL        string    `json:"merges_url"`
	ArchiveURL       string    `json:"archive_url"`
	DownloadsURL     string    `json:"downloads_url"`
	IssuesURL        string    `json:"issues_url"`
	PullsURL         string    `json:"pulls_url"`
	MilestonesURL    string    `json:"milestones_url"`
	NotificationsURL string    `json:"notifications_url"`
	LabelsURL        string    `json:"labels_url"`
	ReleasesURL      string    `json:"releases_url"`
	DeploymentsURL   string    `json:"deployments_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PushedAt         time.Time `json:"pushed_at"`
	GitURL           string    `json:"git_url"`
	SSHURL           string    `json:"ssh_url"`
	CloneURL         string    `json:"clone_url"`
	SvnURL           string    `json:"svn_url"`
	Homepage         string    `json:"homepage"`
	Size             int64     `json:"size"`
	StargazersCount  int64     `json:"stargazers_count"`
	WatchersCount    int64     `json:"watchers_count"`
	Language         string    `json:"language"`
	HasIssues        bool      `json:"has_issues"`
	HasProjects      bool      `json:"has_projects"`
	HasDownloads     bool      `json:"has_downloads"`
	HasWiki          bool      `json:"has_wiki"`
	HasPages         bool      `json:"has_pages"`
	ForksCount       int64     `json:"forks_count"`
	MirrorURL        string    `json:"mirror_url"`
	Archived         bool      `json:"archived"`
	Disabled         bool      `json:"disabled"`
	OpenIssuesCount  int64     `json:"open_issues_count"`
	License          License   `json:"license"`
	Forks            int64     `json:"forks"`
	OpenIssues       int64     `json:"open_issues"`
	Watchers         int64     `json:"watchers"`
	DefaultBranch    string    `json:"default_branch"`
}

type License struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	SpdxID  string `json:"spdx_id"`
	NodeID  string `json:"node_id"`
	HTMLURL string `json:"html_url"`
}

type PullRequest struct {
	URL                string    `json:"url"`
	ID                 int64     `json:"id"`
	NodeID             string    `json:"node_id"`
	HTMLURL            string    `json:"html_url"`
	DiffURL            string    `json:"diff_url"`
	PatchURL           string    `json:"patch_url"`
	IssueURL           string    `json:"issue_url"`
	CommitsURL         string    `json:"commits_url"`
	ReviewCommentsURL  string    `json:"review_comments_url"`
	ReviewCommentURL   string    `json:"review_comment_url"`
	CommentsURL        string    `json:"comments_url"`
	StatusesURL        string    `json:"statuses_url"`
	Number             int64     `json:"number"`
	State              string    `json:"state"`
	Locked             bool      `json:"locked"`
	Title              string    `json:"title"`
	User               User      `json:"user"`
	Body               string    `json:"body"`
	Labels             []Label   `json:"labels"`
	Milestone          Milestone `json:"milestone"`
	ActiveLockReason   string    `json:"active_lock_reason"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ClosedAt           time.Time `json:"closed_at"`
	MergedAt           time.Time `json:"merged_at"`
	MergeCommitSha     string    `json:"merge_commit_sha"`
	Assignee           User      `json:"assignee"`
	Assignees          []User    `json:"assignees"`
	RequestedReviewers []User    `json:"requested_reviewers"`
	RequestedTeams     []Team    `json:"requested_teams"`
	Head               Commit    `json:"head"`
	Base               Commit    `json:"base"`
	Links              struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Issue struct {
			Href string `json:"href"`
		} `json:"issue"`
		Comments struct {
			Href string `json:"href"`
		} `json:"comments"`
		ReviewComments struct {
			Href string `json:"href"`
		} `json:"review_comments"`
		ReviewComment struct {
			Href string `json:"href"`
		} `json:"review_comment"`
		Commits struct {
			Href string `json:"href"`
		} `json:"commits"`
		Statuses struct {
			Href string `json:"href"`
		} `json:"statuses"`
	} `json:"_links"`
	AuthorAssociation   string `json:"author_association"`
	Draft               bool   `json:"draft"`
	Merged              bool   `json:"merged"`
	Mergeable           bool   `json:"mergeable"`
	Rebaseable          bool   `json:"rebaseable"`
	MergeableState      string `json:"mergeable_state"`
	MergedBy            User   `json:"merged_by"`
	Comments            int64  `json:"comments"`
	ReviewComments      int64  `json:"review_comments"`
	MaintainerCanModify bool   `json:"maintainer_can_modify"`
	Commits             int64  `json:"commits"`
	Additions           int64  `json:"additions"`
	Deletions           int64  `json:"deletions"`
	ChangedFiles        int64  `json:"changed_files"`
}

type Label struct {
	ID          int64  `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
}

type Milestone struct {
	URL          string     `json:"url"`
	HTMLURL      string     `json:"html_url"`
	LabelsURL    string     `json:"labels_url"`
	ID           int64      `json:"id"`
	NodeID       string     `json:"node_id"`
	Number       int        `json:"number"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Creator      User       `json:"creator"`
	OpenIssues   int64      `json:"open_issues"`
	ClosedIssues int64      `json:"closed_issues"`
	State        string     `json:"state"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DueOn        *time.Time `json:"due_on"`
	ClosedAt     time.Time  `json:"closed_at"`
}

type User struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Team struct {
	ID              int64  `json:"id"`
	NodeID          string `json:"node_id"`
	URL             string `json:"url"`
	HTMLURL         string `json:"html_url"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Description     string `json:"description"`
	Privacy         string `json:"privacy"`
	Permission      string `json:"permission"`
	MembersURL      string `json:"members_url"`
	RepositoriesURL string `json:"repositories_url"`
}

type Commit struct {
	Label string     `json:"label"`
	Ref   string     `json:"ref"`
	Sha   string     `json:"sha"`
	User  User       `json:"user"`
	Repo  Repository `json:"repo"`
}

type Organization struct {
	Login                                  string    `json:"login"`
	ID                                     int64     `json:"id"`
	NodeID                                 string    `json:"node_id"`
	URL                                    string    `json:"url"`
	ReposURL                               string    `json:"repos_url"`
	EventsURL                              string    `json:"events_url"`
	HooksURL                               string    `json:"hooks_url"`
	IssuesURL                              string    `json:"issues_url"`
	MembersURL                             string    `json:"members_url"`
	PublicMembersURL                       string    `json:"public_members_url"`
	AvatarURL                              string    `json:"avatar_url"`
	Description                            string    `json:"description"`
	Name                                   string    `json:"name"`
	Company                                string    `json:"company"`
	Blog                                   string    `json:"blog"`
	Location                               string    `json:"location"`
	Email                                  string    `json:"email"`
	TwitterUsername                        string    `json:"twitter_username"`
	IsVerified                             bool      `json:"is_verified"`
	HasOrganizationProjects                bool      `json:"has_organization_projects"`
	HasRepositoryProjects                  bool      `json:"has_repository_projects"`
	PublicRepos                            int64     `json:"public_repos"`
	PublicGists                            int64     `json:"public_gists"`
	Followers                              int64     `json:"followers"`
	Following                              int64     `json:"following"`
	HTMLURL                                string    `json:"html_url"`
	CreatedAt                              time.Time `json:"created_at"`
	UpdatedAt                              time.Time `json:"updated_at"`
	Type                                   string    `json:"type"`
	TotalPrivateRepos                      int64     `json:"total_private_repos"`
	OwnedPrivateRepos                      int64     `json:"owned_private_repos"`
	PrivateGists                           int64     `json:"private_gists"`
	DiskUsage                              int64     `json:"disk_usage"`
	Collaborators                          int64     `json:"collaborators"`
	BillingEmail                           string    `json:"billing_email"`
	Plan                                   Plan      `json:"plan"`
	DefaultRepositoryPermission            string    `json:"default_repository_permission"`
	MembersCanCreateRepositories           bool      `json:"members_can_create_repositories"`
	TwoFactorRequirementEnabled            bool      `json:"two_factor_requirement_enabled"`
	MembersAllowedRepositoryCreationType   string    `json:"members_allowed_repository_creation_type"`
	MembersCanCreatePublicRepositories     bool      `json:"members_can_create_public_repositories"`
	MembersCanCreatePrivateRepositories    bool      `json:"members_can_create_private_repositories"`
	MembersCanCreateint64ernalRepositories bool      `json:"members_can_create_int64ernal_repositories"`
	MembersCanCreatePages                  bool      `json:"members_can_create_pages"`
	MembersCanForkPrivateRepositories      bool      `json:"members_can_fork_private_repositories"`
}

type Plan struct {
	Name         string `json:"name"`
	Space        int64  `json:"space"`
	PrivateRepos int64  `json:"private_repos"`
	FilledSeats  int64  `json:"filled_seats"`
	Seats        int64  `json:"seats"`
}
