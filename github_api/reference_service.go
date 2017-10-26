package github_api

import (
	"time"
	"fmt"
)

type Reference struct {
	Ref    string `json:"ref"`
	URL    string `json:"url"`
	Object struct {
		Type string `json:"type"`
		Sha  string `json:"sha"`
		URL  string `json:"url"`
	} `json:"object"`
	CommitDate time.Time
}

type ReferenceService Service

func NewReferenceService(client *GithubClient) *ReferenceService {

	return &ReferenceService{
		client: client,
	}
}

func (rs *ReferenceService) ListReferences() ([]*Reference, error) {
	var allRefs, refs []*Reference

	opts := &ListOpts{}
	for {
		req, err := rs.client.NewRequest("GET", "git/refs/heads", opts)
		if err != nil {
			return nil, err
		}

		_, err = rs.client.Do(req, &refs)
		if err != nil {
			return nil, err
		}
		allRefs = append(allRefs, refs...)

		break

	}

	return allRefs, nil
}

func (rs *ReferenceService) DeleteReference(reference string) (error) {
	u := fmt.Sprintf("git/%s", reference)
	req, err := rs.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = rs.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}