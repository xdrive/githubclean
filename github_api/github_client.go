package github_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
)

const (
	githubBaseUrl = "https://api.github.com/"
	userAgent = "githubclean"
)

type GithubClient struct {
	httpClient *http.Client
	baseUrl    *url.URL
}

func NewGithubClient(githubUser, githubRepository, accessToken string) *GithubClient {
	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: accessToken},
			),
		},
	}

	u := fmt.Sprintf("%srepos/%s/%s/", githubBaseUrl, githubUser, githubRepository)
	baseUrl, _ := url.Parse(u)

	return &GithubClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (c *GithubClient) NewRequest(method, relUrl string, opts *ListOpts) (*http.Request, error) {
	if opts != nil && opts.Page > 0 {
		relUrl = fmt.Sprintf("%s?page=%d", relUrl, opts.Page)
	}
	u, _ := url.Parse(relUrl)
	fullUrl := c.baseUrl.ResolveReference(u)

	req, err := http.NewRequest(method, fullUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func (c *GithubClient) Do(req *http.Request, v interface{}) (*Links, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//dump, _ := httputil.DumpResponse(resp, true)
	//fmt.Printf("%q", dump)

	if err = c.checkResponseStatus(resp, req); err != nil {
		return nil, err
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("get:\n", string(body))

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, err
		}
	}

	return NewLinks(resp), nil
}

func (c *GithubClient) checkResponseStatus(resp *http.Response, req *http.Request) error {
	if req.Method == http.MethodGet && resp.StatusCode != http.StatusOK ||
		req.Method == http.MethodDelete && resp.StatusCode != http.StatusNoContent {

		err := errors.New(fmt.Sprintf("API status code: %s, for url: %s", resp.Status, req.URL))
		return err
	}

	return nil
}

type ListOpts struct {
	Page int
}

type Service struct {
	client *GithubClient
}

type Links struct {
	First int
	Prev  int
	Next  int
	Last  int
}

func NewLinks(resp *http.Response) *Links {
	first, prev, next, last := parseHeaderLinks(resp)
	l := &Links{
		First: first,
		Prev:  prev,
		Next:  next,
		Last:  last,
	}

	return l
}

func parseHeaderLinks(resp *http.Response) (first int, prev int, next int, last int) {
	if links := resp.Header.Get("Link"); links != "" {
		for _, link := range strings.Split(links, ", ") {
			linkParts := strings.Split(link, "; ")
			if len(linkParts) != 2 {
				continue
			}

			u, err := url.Parse(linkParts[0][1 : len(linkParts[0])-1])
			if err != nil {
				continue
			}

			pageS := u.Query().Get("page")
			page, err := strconv.Atoi(pageS)
			if err != nil {
				continue
			}

			switch linkParts[1] {
			case `rel="next"`:
				next = page
			case `rel="last"`:
				last = page
			case `rel="first"`:
				first = page
			case `rel="prev"`:
				prev = page
			}

		}
	}

	return
}
