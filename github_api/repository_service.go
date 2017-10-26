package github_api

import (
	"fmt"
	"time"
)

type RepositoryCommit struct {
	URL         string `json:"url"`
	Sha         string `json:"sha"`
	Commit      struct {
		URL    string `json:"url"`
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

type Branch struct {
	Name   string `json:"name"`
	Commit struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	Protected     bool   `json:"protected"`
	ProtectionURL string `json:"protection_url"`
}


type RepositoryService Service

func NewRepositoryService(client *GithubClient) *RepositoryService {

	return &RepositoryService{
		client: client,
	}
}

func (rs *RepositoryService) GetCommit(sha1 string) (*RepositoryCommit, error) {
	u := fmt.Sprintf("commits/%s", sha1)
	req, err := rs.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var repositoryCommit RepositoryCommit

	_, err = rs.client.Do(req, &repositoryCommit)
	if err != nil {
		return nil, err
	}

	//if err := json.NewDecoder(resp.Body).Decode(&repositoryCommit); err != nil {
	//	return nil, err
	//}

	return &repositoryCommit, nil
}

func (rs *RepositoryService) ListBranches() ([]*Branch, error) {
	var i int
	var allBranches, branches []*Branch

	opts := &ListOpts{}
	for {
		req, err := rs.client.NewRequest("GET", "branches", opts)
		if err != nil {
			return nil, err
		}

		links, err := rs.client.Do(req, &branches)
		if err != nil {
			return nil, err
		}
		allBranches = append(allBranches, branches...)
		branches = nil

		if links.Next == 0 {
			break
		}

		opts.Page = links.Next

		// todo remove. added to not request too much branches
		i++
		if i > 1 {
			break
		}
	}


	//if links, ok := resp.Header["Link"]; ok {
	//	fmt.Println(links)
	//}

	//dump, _ := httputil.DumpResponse(resp, true)
	//fmt.Printf("%q", dump)

	//if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
	//	return nil, err
	//}

	return allBranches, nil
}