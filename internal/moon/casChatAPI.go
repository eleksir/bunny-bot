package moon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hjson/hjson-go"
)

// CasCheckID check given telegram user id vs CAS spammers base.
// https://api.cas.chat/check?user_id=821871410
func CasCheckID(id int64) (bool, error) {
	var (
		result bool
		err    error
		url    = fmt.Sprintf("https://api.cas.chat/check?user_id=%d", id)
		c      = &http.Client{
			Timeout: 15 * time.Second,
		}
		banned    CasTrue
		notBanned CasFalse
	)

	resp, err := c.Get(url) //nolint: bodyclose, noctx

	if err != nil {
		return result, fmt.Errorf("unable to query api.cas.chat, %w", err)
	}

	// Looks like we should clse body as it it recommends in docs: https://pkg.go.dev/net/http .
	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			Log.Errorf("Unable to close response body for request to %s: %s", url, err)
		}
	}(resp.Body)

	// Client should not return with http status code other than 200.
	if resp.StatusCode > 200 {
		return result, fmt.Errorf( //nolint: err113
			"api.cas.chat returns non-200 http status code: %d %s",
			resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return result, fmt.Errorf("unable to read response from api.cas.chat, %w", err)
	}

	// looks like go native json decoder is pretty shitty, so we'll take hjson as less restrictive parser, parse json
	// encode it back to struct, then parse struct as encode/json

	var tmp map[string]any
	err = hjson.Unmarshal(body, &tmp)

	// Не удалось распарсить.
	if err != nil {
		return result, fmt.Errorf("unable to parse response from api.cas.chat via hjson: %w", err)
	}

	tmpjson, err := json.Marshal(tmp)

	// Не удалось преобразовать map-ку в json.
	if err != nil {
		return result, fmt.Errorf(
			"unable to parse back to json response from api.cas.chat parsed by hjson: %w", err)
	}

	if err := json.Unmarshal(tmpjson, &notBanned); err == nil {
		// if ok field in json == false, then user is not in ban. (records about id not found)
		if !notBanned.Ok {
			return false, err
		}
	} // At this point we just ignore err != nil, because it *can be* other json type.

	if err := json.Unmarshal(body, &banned); err == nil {
		// if ok field in json == true, then user is banned. (records about id are found)
		if banned.Ok {
			Log.Debugf(
				"api.cas.chat reported user with id=%d as banned: offences=%d ban_time=%s",
				id,
				banned.Result.Offenses,
				banned.Result.TimeAdded,
			)

			return true, err
		}

		return result, fmt.Errorf("api.cas.chat returns strange structure: %s", string(body)) //nolint: err113
	} // At this point we can't ignore err != nil because we do not know (yet?) other forms of json reply from cas api.

	return result, fmt.Errorf("unable to parse response from api.cas.chat: %s", string(body)) //nolint: err113
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
