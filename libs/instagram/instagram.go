package instagram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"io/ioutil"

	"github.com/pkg/errors"
)

const (
	apiUrl      = "https://www.instagram.com/"
	mediaPath   = "p/"
	httpTimeout = time.Second * 3
)

var (
	MediaForbidden = errors.New("it's a private account or the media doesn't exist")
)

var (
	httpClient = http.Client{Timeout: httpTimeout}
)

type Media struct {
	Owner   string
	Url     string
	IsPhoto bool
}

type mediaResponse struct {
	EntryData struct {
		PostPage []struct {
			GraphQL struct {
				Media struct {
					Owner struct {
						Username string
					}
					DisplayResources []struct {
						Src string
					} `json:"display_resources"`
					IsVideo *bool `json:"is_video"`
					Sidecar struct {
						Edges []struct {
							Node struct {
								IsVideo *bool `json:"is_video"`
							}
						}
					} `json:"edge_sidecar_to_children"`
				} `json:"shortcode_media"`
			} `json:"graphql"`
		}
	} `json:"entry_data"`
}

func GetMediaInfo(ctx context.Context, mediaCode string) (*Media, error) {
	mediaUrl := apiUrl + mediaPath + mediaCode
	resp, err := httpClient.Get(mediaUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, MediaForbidden
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("media info response unexpected status: %d", resp.StatusCode)
	}
	responseData := &mediaResponse{}
	err = extractData(resp, responseData)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse media response")
	}
	err = validateMediaResponse(responseData)
	if err != nil {
		return nil, errors.Wrap(err, "invalid media response")
	}
	mediaData := responseData.EntryData.PostPage[0].GraphQL.Media
	username := mediaData.Owner.Username
	displayUrl := mediaData.DisplayResources[0].Src
	isPhoto := true
	if *mediaData.IsVideo {
		isPhoto = false
	} else {
		edges := mediaData.Sidecar.Edges
		if len(edges) != 0 {
			if *edges[0].Node.IsVideo {
				isPhoto = false
			}
		}
	}

	return &Media{Owner: username, Url: displayUrl, IsPhoto: isPhoto}, nil
}

func validateMediaResponse(data *mediaResponse) error {
	if len(data.EntryData.PostPage) != 1 {
		return errors.Errorf("expected only one post, got: %d", len(data.EntryData.PostPage))
	}
	media := data.EntryData.PostPage[0].GraphQL.Media
	if media.Owner.Username == "" {
		return errors.New("no username")
	}
	if len(media.DisplayResources) == 0 {
		return errors.New("no display resources")
	}
	if media.DisplayResources[0].Src == "" {
		return errors.New("no display resource src")
	}
	if media.IsVideo == nil {
		return errors.New("no is video")
	}
	if len(media.Sidecar.Edges) != 0 && media.Sidecar.Edges[0].Node.IsVideo == nil {
		return errors.New("no sidecar is video")
	}
	return nil
}

type userResponse struct {
	EntryData struct {
		ProfilePage []struct {
			GraphQL struct {
				User struct {
					FollowedBy struct {
						Count *int
					} `json:"edge_followed_by"`
				}
			} `json:"graphql"`
		}
	} `json:"entry_data"`
}

func GetFollowersNumber(ctx context.Context, username string) (int, error) {
	profileUrl := apiUrl + username
	resp, err := httpClient.Get(profileUrl)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errors.Errorf("user response unexpected status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	responseData := &userResponse{}
	err = extractData(resp, responseData)
	if err != nil {
		return 0, errors.Wrap(err, "can't parse user response")
	}
	err = validateUserResponse(responseData)
	if err != nil {
		return 0, errors.Wrap(err, "invalid user response")
	}
	return *responseData.EntryData.ProfilePage[0].GraphQL.User.FollowedBy.Count, nil
}

func validateUserResponse(data *userResponse) error {
	if len(data.EntryData.ProfilePage) != 1 {
		return errors.Errorf("expected only one profile, got: %d", len(data.EntryData.ProfilePage))
	}
	if data.EntryData.ProfilePage[0].GraphQL.User.FollowedBy.Count == nil {
		return errors.New("no followers info")
	}
	return nil
}

func extractData(resp *http.Response, destination interface{}) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	text := string(b)
	firstMarker := "window._sharedData = "
	secondMarker := ";</script>"
	dataStart := strings.Index(text, firstMarker) + len(firstMarker)
	dataEnd := strings.Index(text, secondMarker)
	data := text[dataStart:dataEnd]
	return json.Unmarshal([]byte(data), destination)
}

func BuildProfileUrl(username string) string {
	return apiUrl + username + "/"
}

func ExtractUsername(profileUrl string) (string, error) {
	username, err := extractLastPathPart(profileUrl)
	return username, errors.Wrap(err, "cant extract username")
}

func ExtractMediaCode(mediaUrl string) (string, error) {
	mediaCode, err := extractLastPathPart(mediaUrl)
	return mediaCode, err
}

func extractLastPathPart(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	urlPath := strings.TrimSuffix(u.Path, "/")
	_, lastPart := path.Split(urlPath)
	if lastPart == "" {
		return "", errors.New("url has empty path")
	}
	return lastPart, nil
}
