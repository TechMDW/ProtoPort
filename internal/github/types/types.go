package types

type GithubContentApi struct {
	Type        string                `json:"type,omitempty"`
	Encoding    string                `json:"encoding,omitempty"`
	Size        int                   `json:"size,omitempty"`
	Name        string                `json:"name,omitempty"`
	Path        string                `json:"path,omitempty"`
	Content     string                `json:"content,omitempty"`
	Sha         string                `json:"sha,omitempty"`
	URL         string                `json:"url,omitempty"`
	GitURL      string                `json:"git_url,omitempty"`
	HTMLURL     string                `json:"html_url,omitempty"`
	DownloadURL string                `json:"download_url,omitempty"`
	Links       GithubContentLinksApi `json:"_links,omitempty"`
}
type GithubContentLinksApi struct {
	Git  string `json:"git,omitempty"`
	Self string `json:"self,omitempty"`
	HTML string `json:"html,omitempty"`
}
