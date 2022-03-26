package cli

import (
	"net/http"
	"path"
	"strconv"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-auth/auth"
	"github.com/filecoin-project/venus-auth/config"
	"github.com/filecoin-project/venus-auth/core"
	"github.com/filecoin-project/venus-auth/errcode"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

type localClient struct {
	cli *resty.Client
}

// nolint
func GetCli(ctx *cli.Context) (*localClient, error) {
	p, err := homedir.Expand(ctx.String("repo"))
	if err != nil {
		return nil, xerrors.Errorf("could not expand home dir (repo): %w", err)
	}
	cnfPath, err := homedir.Expand(ctx.String("config"))
	if err != nil {
		return nil, xerrors.Errorf("could not expand home dir (config): %w", err)
	}
	if len(cnfPath) == 0 {
		cnfPath = path.Join(p, "config.toml")
	}
	cnf, err := config.DecodeConfig(cnfPath)
	if err != nil {
		return nil, xerrors.Errorf("failed to decode config err: %w", err)
	}
	c, err := newClient(cnf.Port)
	return c, err
}

func newClient(port string) (*localClient, error) {
	client := resty.New().
		SetHostURL("http://localhost:"+port).
		SetHeader("Accept", "application/json")
	return &localClient{cli: client}, nil
}

func (lc *localClient) GenerateToken(name, perm, extra string) (string, error) {
	resp, err := lc.cli.R().SetBody(auth.GenTokenRequest{
		Name:  name,
		Perm:  perm,
		Extra: extra,
	}).SetResult(&auth.GenTokenResponse{}).SetError(&errcode.ErrMsg{}).Post("/genToken")
	if err != nil {
		return core.EmptyString, err
	}
	if resp.StatusCode() == http.StatusOK {
		res := resp.Result().(*auth.GenTokenResponse)
		return res.Token, nil
	}
	return core.EmptyString, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) Tokens(skip, limit int64) (auth.GetTokensResponse, error) {
	resp, err := lc.cli.R().SetQueryParams(map[string]string{
		"skip":  strconv.FormatInt(skip, 10),
		"limit": strconv.FormatInt(limit, 10),
	}).SetResult(&auth.GetTokensResponse{}).SetError(&errcode.ErrMsg{}).Get("/tokens")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		return *(resp.Result().(*auth.GetTokensResponse)), nil
	}
	return nil, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) RemoveToken(token string) error {
	resp, err := lc.cli.R().SetBody(auth.RemoveTokenRequest{
		Token: token,
	}).SetError(&errcode.ErrMsg{}).Delete("/token")
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	}
	return resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) CreateUser(req *auth.CreateUserRequest) (*auth.CreateUserResponse, error) {
	resp, err := lc.cli.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&auth.CreateUserResponse{}).
		SetError(&errcode.ErrMsg{}).
		Put("/user/new")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		return resp.Result().(*auth.CreateUserResponse), nil
	}
	return nil, resp.Error().(*errcode.ErrMsg).Err()
}

// UpdateUser
func (lc *localClient) UpdateUser(req *auth.UpdateUserRequest) error {
	resp, err := lc.cli.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).SetError(&errcode.ErrMsg{}).Post("/user/update")
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	}
	return resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) ListUsers(req *auth.ListUsersRequest) (auth.ListUsersResponse, error) {
	resp, err := lc.cli.R().SetQueryParams(map[string]string{
		"skip":       strconv.FormatInt(req.Skip, 10),
		"limit":      strconv.FormatInt(req.Limit, 10),
		"sourceType": strconv.Itoa(req.SourceType),
		"state":      strconv.Itoa(req.State),
		"keySum":     strconv.Itoa(req.KeySum),
	}).SetResult(&auth.ListUsersResponse{}).SetError(&errcode.ErrMsg{}).Get("/user/list")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		return *(resp.Result().(*auth.ListUsersResponse)), nil
	}
	return nil, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) GetUser(req *auth.GetUserRequest) (*auth.OutputUser, error) {
	resp, err := lc.cli.R().SetQueryParams(map[string]string{
		"name": req.Name,
	}).SetResult(&auth.OutputUser{}).SetError(&errcode.ErrMsg{}).Get("/user")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		return resp.Result().(*auth.OutputUser), nil
	}
	return nil, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) GetUserByMiner(req *auth.GetUserByMinerRequest) (*auth.OutputUser, error) {
	resp, err := lc.cli.R().SetQueryParams(map[string]string{
		"miner": req.Miner,
	}).SetResult(&auth.OutputUser{}).SetError(&errcode.ErrMsg{}).Get("/miner")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		return resp.Result().(*auth.OutputUser), nil
	}
	return nil, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) HasMiner(req *auth.HasMinerRequest) (bool, error) {
	var has bool
	resp, err := lc.cli.R().SetQueryParams(map[string]string{
		"miner": req.Miner,
	}).SetResult(&has).SetError(&errcode.ErrMsg{}).Get("/miner/has-miner")
	if err != nil {
		return false, err
	}

	if resp.StatusCode() == http.StatusOK {
		return *resp.Result().(*bool), nil
	}
	return false, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) GetUserRateLimit(name, id string) (auth.GetUserRateLimitResponse, error) {
	var res auth.GetUserRateLimitResponse
	resp, err := lc.cli.R().SetQueryParams(
		map[string]string{"name": name, "id": id}).
		SetResult(&res).
		SetError(&errcode.ErrMsg{}).
		Get("/user/ratelimit")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		return *(resp.Result().(*auth.GetUserRateLimitResponse)), nil
	}
	return nil, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) UpsertUserRateLimit(req *auth.UpsertUserRateLimitReq) (string, error) {
	var res string
	resp, err := lc.cli.R().SetBody(req).SetResult(&res).SetError(&errcode.ErrMsg{}).Post("/user/ratelimit/upsert")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() == http.StatusOK {
		return *(resp.Result().(*string)), nil
	}
	return "", resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) DelUserRateLimit(req *auth.DelUserRateLimitReq) (string, error) {
	var id string
	resp, err := lc.cli.R().SetBody(req).SetResult(&id).SetError(&errcode.ErrMsg{}).Post("/user/ratelimit/del")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() == http.StatusOK {
		return *(resp.Result().(*string)), nil
	}
	return "", resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) UpsertMiner(user, miner string) (bool, error) {
	if _, err := address.NewFromString(miner); err != nil {
		return false, xerrors.Errorf("invalid miner address:%s", miner)
	}
	var isCreate bool
	resp, err := lc.cli.R().SetBody(&auth.UpsertMinerReq{Miner: miner, User: user}).
		SetResult(&isCreate).SetError(&errcode.ErrMsg{}).Post("/miner/add-miner")
	if err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusOK {
		return *(resp.Result().(*bool)), nil
	}
	return false, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) ListMiners(user string) (auth.ListMinerResp, error) {
	var res auth.ListMinerResp
	resp, err := lc.cli.R().SetQueryParams(map[string]string{"user": user}).
		SetResult(&res).SetError(&errcode.ErrMsg{}).Get("/miner/list-by-user")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		return *(resp.Result().(*auth.ListMinerResp)), nil
	}
	return nil, resp.Error().(*errcode.ErrMsg).Err()
}

func (lc *localClient) DelMiner(miner string) (bool, error) {
	if _, err := address.NewFromString(miner); err != nil {
		return false, xerrors.Errorf("invalid miner address:%s", miner)
	}

	var res bool
	resp, err := lc.cli.R().SetBody(auth.DelMinerReq{Miner: miner}).
		SetResult(&res).SetError(&errcode.ErrMsg{}).Post("/miner/del")
	if err != nil {
		return res, err
	}
	if resp.StatusCode() == http.StatusOK {
		return *(resp.Result().(*bool)), nil
	}
	return res, resp.Error().(*errcode.ErrMsg).Err()
}
