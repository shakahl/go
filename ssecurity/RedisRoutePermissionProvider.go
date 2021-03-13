package ssecurity

import (
	"encoding/json"

	"github.com/go-redis/redis/v7"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/sproto"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/u"
)

type RedisRoutePermissionProvider struct {
	redis         redis.Cmdable
	RouteKey      string
	PermissionKey string
}

func NewRedisRoutePermissionProvider(routeKey, permissionKey string, config *sredis.RedisConfig) IRoutePermissionProvider {
	if routeKey == "" {
		log.Fatal("routeKey cannot be empty")
	}

	if permissionKey == "" {
		log.Fatal("permissionKey key cannot be empty")
	}

	r := new(RedisRoutePermissionProvider)

	r.redis = sredis.NewClient(config)

	r.RouteKey = routeKey
	r.PermissionKey = permissionKey

	return r
}

// *******************************************************************************************************************************
// Route
func (x *RedisRoutePermissionProvider) CreateRoute(in *sproto.RouteDTO) error {
	j, err := json.Marshal(in)
	if err != nil {
		return err
	}

	cmd := x.redis.HSet(x.RouteKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRoutePermissionProvider) GetRoute(id string) (*sproto.RouteDTO, error) {
	cmd := x.redis.HGet(x.RouteKey, id)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}

	r := new(sproto.RouteDTO)
	j, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(j), r)
	u.LogError(err)
	return r, err
}
func (x *RedisRoutePermissionProvider) UpdateRoute(in *sproto.RouteDTO) error {
	j, err := json.Marshal(in)
	if err != nil {
		return err
	}

	cmd := x.redis.HSet(x.RouteKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRoutePermissionProvider) RemoveRoute(id string) error {
	cmd := x.redis.HDel(x.RouteKey, id)
	return cmd.Err()
}
func (x *RedisRoutePermissionProvider) GetRoutes() (map[string]*sproto.RouteDTO, error) {
	cmd := x.redis.HGetAll(x.RouteKey)
	r, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	m := make(map[string]*sproto.RouteDTO, len(r))
	for key, value := range r {
		dto := new(sproto.RouteDTO)
		err = json.Unmarshal([]byte(value), dto)
		if !u.LogError(err) {
			m[key] = dto
		}
	}
	return m, err
}

// *******************************************************************************************************************************
// Permission
func (x *RedisRoutePermissionProvider) CreatePermission(in *sproto.PermissionDTO) error {
	j, err := json.Marshal(in)
	if err != nil {
		return err
	}

	cmd := x.redis.HSet(x.PermissionKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRoutePermissionProvider) GetPermission(id string) (*sproto.PermissionDTO, error) {
	cmd := x.redis.HGet(x.PermissionKey, id)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}

	r := new(sproto.PermissionDTO)
	j, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(j), r)
	u.LogError(err)
	return r, err
}
func (x *RedisRoutePermissionProvider) UpdatePermission(in *sproto.PermissionDTO) error {
	j, err := json.Marshal(in)
	if err != nil {
		return err
	}

	cmd := x.redis.HSet(x.PermissionKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRoutePermissionProvider) RemovePermission(id string) error {
	cmd := x.redis.HDel(x.PermissionKey, id)
	return cmd.Err()
}
func (x *RedisRoutePermissionProvider) GetPermissions() (map[string]*sproto.PermissionDTO, error) {
	cmd := x.redis.HGetAll(x.PermissionKey)
	r, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	m := make(map[string]*sproto.PermissionDTO, len(r))
	for key, value := range r {
		dto := new(sproto.PermissionDTO)
		err = json.Unmarshal([]byte(value), dto)
		if err == nil {
			m[key] = dto
		} else {
			log.Errorf("%s, %v", value, err)
		}
	}
	return m, err
}