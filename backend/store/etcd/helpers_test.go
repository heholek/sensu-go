package etcd

import (
	"context"
	"testing"

	"github.com/sensu/sensu-go/backend/store"
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	testWithEtcd(t, func(store store.Store) {
		etcd, _ := store.(*etcdStore)

		// Create a new org "acme" and its environments "default" & "dev"
		store.UpdateOrganization(context.Background(), &types.Organization{
			Name: "acme",
		})
		store.UpdateEnvironment(context.Background(),
			"acme",
			&types.Environment{
				Name: "default",
			})
		store.UpdateEnvironment(context.Background(),
			"acme",
			&types.Environment{
				Name: "dev",
			})

		// Create /checks/default/default/check1
		check1 := types.FixtureCheckConfig("check1")
		// ctx := context.WithValue(context.Background(), types.OrganizationKey, check1.Organization)
		if err := store.UpdateCheckConfig(context.Background(), check1); err != nil {
			assert.FailNow(t, err.Error())
		}

		// Create /checks/acme/default/check2
		check2 := types.FixtureCheckConfig("check2")
		check2.Organization = "acme"
		// ctx = context.WithValue(context.Background(), types.OrganizationKey, check2.Organization)
		if err := store.UpdateCheckConfig(context.Background(), check2); err != nil {
			assert.FailNow(t, err.Error())
		}

		// Create /checks/acme/dev/check3
		check3 := types.FixtureCheckConfig("check3")
		check3.Organization = "acme"
		check3.Environment = "dev"
		// ctx = context.WithValue(context.Background(), types.OrganizationKey, check3.Organization)
		if err := store.UpdateCheckConfig(context.Background(), check3); err != nil {
			assert.FailNow(t, err.Error())
		}

		// Mock a context that put ourselves in the default/default environment
		ctx := context.WithValue(context.Background(), types.OrganizationKey, "default")
		ctx = context.WithValue(ctx, types.EnvironmentKey, "default")

		// We only have a single result given our current org & env
		resp, err := query(ctx, etcd, getCheckConfigsPath)
		assert.NoError(t, err)
		assert.Len(t, resp.Kvs, 1)

		// Mock a context to query across every single organization
		ctx = context.WithValue(ctx, types.OrganizationKey, "*")

		// We now have two result given our "wildcard" org
		resp, err = query(ctx, etcd, getCheckConfigsPath)
		assert.NoError(t, err)
		assert.Len(t, resp.Kvs, 2)

		// Mock a context to query across every single environment of the acme org
		ctx = context.WithValue(ctx, types.OrganizationKey, "acme")
		ctx = context.WithValue(ctx, types.EnvironmentKey, "*")

		// We now have two result given our "wildcard" env
		resp, err = query(ctx, etcd, getCheckConfigsPath)
		assert.NoError(t, err)
		assert.Len(t, resp.Kvs, 2)

		// Mock a context to query across every single organization and environment
		ctx = context.WithValue(ctx, types.OrganizationKey, "*")
		ctx = context.WithValue(ctx, types.EnvironmentKey, "*")

		// We now have two result given our "wildcard" org
		resp, err = query(ctx, etcd, getCheckConfigsPath)
		assert.NoError(t, err)
		assert.Len(t, resp.Kvs, 3)
	})
}
