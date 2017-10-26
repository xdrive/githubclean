package github_api

type PullRequest struct {
	ID                int    `json:"id"`
	URL               string `json:"url"`
	Head      struct {
		Ref   string `json:"ref"`
		Sha   string `json:"sha"`
	} `json:"head"`
}

type PullRequestService Service

func NewPullRequestService(client *GithubClient) *PullRequestService {

	return &PullRequestService{
		client: client,
	}
}

func (rs *PullRequestService) ListPullRequests() ([]*PullRequest, error) {
	var allPulls, pulls []*PullRequest

	opts := &ListOpts{}
	for {
		req, err := rs.client.NewRequest("GET", "pulls", opts)
		if err != nil {
			return nil, err
		}

		_, err = rs.client.Do(req, &pulls)
		if err != nil {
			return nil, err
		}

		allPulls = append(allPulls, pulls...)

		break
	}

	return allPulls, nil
}