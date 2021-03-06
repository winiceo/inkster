package api

import (
	"context"

	"github.com/go-pg/pg"
	gql "github.com/graph-gophers/graphql-go"
)

type pageArgs struct {
	ID gql.ID
}

func (res *Resolver) Page(
	ctx context.Context,
	args pageArgs,
) (*pageResolver, error) {
	localID, err := fromGlobalID("page", string(args.ID))
	if err != nil {
		return nil, err
	}
	result, err := res.dataSource.GetPage(localID)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err

	}

	if !result.IsPublished && !checkPermission(ctx) {
		return nil, errNoPermissions
	}

	return &pageResolver{
		dataSource: res.dataSource,
		data:       &result,
	}, nil
}

type pagesArgs struct{}

func (res *Resolver) Pages(
	ctx context.Context,
) (*[]*pageResolver, error) {
	pages, err := res.dataSource.GetPages()
	if err != nil {
		return nil, err
	}

	resolvers := make([]*pageResolver, len(pages))
	for i, page := range pages {
		resolvers[i] = &pageResolver{
			dataSource: res.dataSource,
			data:       &page,
		}
	}

	return &resolvers, nil
}
