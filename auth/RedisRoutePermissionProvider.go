package auth

import (
	"github.com/go-redis/redis"
	log "github.com/kataras/golog"
	"github.com/syncfuture/go/goredis"
	"github.com/syncfuture/go/json"
	"github.com/syncfuture/go/sproto"
	u "github.com/syncfuture/go/util"
)

const (
	route_key      = "ecp:ROUTES:"
	permission_key = "ecp:PERMISSIONS"
)

func NewRedisRoutePermissionProvider(projectName string, clusterEnabled bool, password string, addrs ...string) IRoutePermissionProvider {
	r := new(RedisRoutePermissionProvider)

	r.redis = goredis.NewClient(clusterEnabled, password, addrs...)
	r.RouteKey = route_key + projectName
	r.PermissionKey = permission_key

	return r
}

type RedisRoutePermissionProvider struct {
	redis         redis.Cmdable
	RouteKey      string
	PermissionKey string
}

// *******************************************************************************************************************************
// Route
func (x *RedisRoutePermissionProvider) CreateRoute(in *sproto.RouteDTO) error {
	j, err := json.Serialize(in)
	if u.LogError(err) {
		return err
	}

	cmd := x.redis.HSet(x.RouteKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRoutePermissionProvider) GetRoute(id string) (*sproto.RouteDTO, error) {
	cmd := x.redis.HGet(x.RouteKey, id)
	err := cmd.Err()
	if u.LogError(err) {
		return nil, err
	}

	r := new(sproto.RouteDTO)
	j, err := cmd.Result()
	if u.LogError(err) {
		return nil, err
	}

	err = json.Deserialize(j, r)
	u.LogError(err)
	return r, err
}
func (x *RedisRoutePermissionProvider) UpdateRoute(in *sproto.RouteDTO) error {
	j, err := json.Serialize(in)
	if u.LogError(err) {
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
	if u.LogError(err) {
		return nil, err
	}

	m := make(map[string]*sproto.RouteDTO, len(r))
	for key, value := range r {
		dto := new(sproto.RouteDTO)
		err = json.Deserialize(value, dto)
		if !u.LogError(err) {
			m[key] = dto
		}
	}
	return m, err
}

// *******************************************************************************************************************************
// Permission
func (x *RedisRoutePermissionProvider) CreatePermission(in *sproto.PermissionDTO) error {
	j, err := json.Serialize(in)
	if u.LogError(err) {
		return err
	}

	cmd := x.redis.HSet(x.PermissionKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRoutePermissionProvider) GetPermission(id string) (*sproto.PermissionDTO, error) {
	cmd := x.redis.HGet(x.PermissionKey, id)
	err := cmd.Err()
	if u.LogError(err) {
		return nil, err
	}

	r := new(sproto.PermissionDTO)
	j, err := cmd.Result()
	if u.LogError(err) {
		return nil, err
	}

	err = json.Deserialize(j, r)
	u.LogError(err)
	return r, err
}
func (x *RedisRoutePermissionProvider) UpdatePermission(in *sproto.PermissionDTO) error {
	j, err := json.Serialize(in)
	if u.LogError(err) {
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
	if u.LogError(err) {
		return nil, err
	}

	m := make(map[string]*sproto.PermissionDTO, len(r))
	for key, value := range r {
		dto := new(sproto.PermissionDTO)
		err = json.Deserialize(value, dto)
		if err == nil {
			m[key] = dto
		} else {
			log.Errorf("%s, %v", value, err)
		}
	}
	return m, err
}
