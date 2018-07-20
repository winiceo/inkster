package api

import (
	"github.com/dominik-zeglen/ecoknow/core"
	"github.com/globalsign/mgo/bson"
	gql "github.com/graph-gophers/graphql-go"
)

type pageCreateResult struct {
	userErrors *[]userError
	page       core.Page
}
type pageRemoveResult struct {
	removedObjectID gql.ID
}
type pageFieldOperationResult struct {
	userErrors *[]userError
	pageID     bson.ObjectId
	page       *core.Page
}

// Type resolvers
type pageResolver struct {
	dataSource core.Adapter
	data       *core.Page
}
type pageFieldResolver struct {
	dataSource core.Adapter
	data       *core.PageField
}
type pageCreateResultResolver struct {
	dataSource core.Adapter
	data       pageCreateResult
}
type pageRemoveResultResolver struct {
	dataSource core.Adapter
	data       pageRemoveResult
}
type pageFieldOperationResultResolver struct {
	dataSource core.Adapter
	data       pageFieldOperationResult
}

func (res *pageRemoveResultResolver) RemovedObjectID() gql.ID {
	return res.data.removedObjectID
}

func (res *pageFieldOperationResultResolver) Page() (*pageResolver, error) {
	result, err := res.dataSource.GetPage(res.data.pageID)
	if err != nil {
		return nil, err
	}
	return &pageResolver{
		dataSource: res.dataSource,
		data:       &result,
	}, nil
}

func (res *pageFieldOperationResultResolver) UserErrors() *[]*userErrorResolver {
	var resolverList []*userErrorResolver
	if res.data.userErrors == nil {
		return nil
	}
	userErrors := *res.data.userErrors
	for i := range userErrors {
		resolverList = append(
			resolverList,
			&userErrorResolver{
				data: userErrors[i],
			},
		)
	}
	return &resolverList
}
func (res *pageCreateResultResolver) Page() *pageResolver {
	return &pageResolver{
		dataSource: res.dataSource,
		data:       &res.data.page,
	}
}

func (res *pageCreateResultResolver) UserErrors() *[]*userErrorResolver {
	var resolverList []*userErrorResolver
	if res.data.userErrors == nil {
		return nil
	}
	userErrors := *res.data.userErrors
	for i := range userErrors {
		resolverList = append(
			resolverList,
			&userErrorResolver{
				data: userErrors[i],
			},
		)
	}
	return &resolverList
}
func (res *pageFieldResolver) Name() string {
	return res.data.Name
}

func (res *pageFieldResolver) Type() string {
	return res.data.Type
}

func (res *pageFieldResolver) Value() *string {
	return &res.data.Value
}

func (res *pageResolver) ID() gql.ID {
	globalID := toGlobalID("page", res.data.ID)
	return gql.ID(globalID)
}

func (res *pageResolver) Name() string {
	return res.data.Name
}

func (res *pageResolver) Slug() string {
	return res.data.Slug
}

func (res *pageResolver) Fields() *[]*pageFieldResolver {
	var resolverList []*pageFieldResolver
	fields := res.data.Fields
	if fields == nil {
		return nil
	}
	for i := range fields {
		resolverList = append(
			resolverList,
			&pageFieldResolver{
				dataSource: res.dataSource,
				data:       &fields[i],
			},
		)
	}
	return &resolverList
}

func (res *pageResolver) Parent() (*directoryResolver, error) {
	parent, err := res.dataSource.GetDirectory(res.data.ParentID)
	if err != nil {
		return nil, err
	}
	return &directoryResolver{
		dataSource: res.dataSource,
		data:       &parent,
	}, nil
}

type createPageArgs struct {
	Input struct {
		Name     string
		ParentID string
		Fields   *[]*struct {
			Name  string
			Type  string
			Value string
		}
	}
}

func (res *Resolver) CreatePage(args createPageArgs) (*pageCreateResultResolver, error) {
	if args.Input.Name == "" {
		return &pageCreateResultResolver{
			dataSource: res.dataSource,
			data: pageCreateResult{
				userErrors: &[]userError{
					userError{
						field:   "name",
						message: errNoEmpty("name").Error(),
					},
				},
			},
		}, nil
	}
	if args.Input.ParentID == "" {
		return &pageCreateResultResolver{
			dataSource: res.dataSource,
			data: pageCreateResult{
				userErrors: &[]userError{
					userError{
						field:   "parentId",
						message: errNoEmpty("parentId").Error(),
					},
				},
			},
		}, nil
	}
	localID, err := fromGlobalID("directory", args.Input.ParentID)
	if err != nil {
		return nil, err
	}
	page := core.Page{
		Name:     args.Input.Name,
		ParentID: bson.ObjectId(localID),
	}
	if args.Input.Fields != nil {
		fields := *args.Input.Fields
		page.Fields = make([]core.PageField, len(fields))
		for i := range fields {
			page.Fields[i] = core.PageField{
				Name:  fields[i].Name,
				Type:  fields[i].Type,
				Value: fields[i].Value,
			}
		}
	}
	result, err := res.dataSource.AddPage(page)
	if err != nil {
		return nil, err
	}
	return &pageCreateResultResolver{
		dataSource: res.dataSource,
		data: pageCreateResult{
			page: result,
		},
	}, nil
}

type createPageFieldArgs struct {
	ID    gql.ID
	Input struct {
		Name  string
		Type  string
		Value string
	}
}

func (res *Resolver) CreatePageField(args createPageFieldArgs) (*pageFieldOperationResultResolver, error) {
	localID, err := fromGlobalID("page", string(args.ID))
	if err != nil {
		return nil, err
	}
	if len(args.Input.Name) == 0 {
		return &pageFieldOperationResultResolver{
			dataSource: res.dataSource,
			data: pageFieldOperationResult{
				userErrors: &[]userError{
					userError{
						field:   "name",
						message: errNoEmpty("name").Error(),
					},
				},
			},
		}, nil
	}
	if len(args.Input.Type) == 0 {
		return &pageFieldOperationResultResolver{
			dataSource: res.dataSource,
			data: pageFieldOperationResult{
				userErrors: &[]userError{
					userError{
						field:   "type",
						message: errNoEmpty("type").Error(),
					},
				},
			},
		}, nil
	}
	field := core.PageField{
		Name:  args.Input.Name,
		Type:  args.Input.Type,
		Value: args.Input.Value,
	}
	err = res.dataSource.AddPageField(localID, field)
	if err != nil {
		return nil, err
	}
	return &pageFieldOperationResultResolver{
		dataSource: res.dataSource,
		data: pageFieldOperationResult{
			pageID: localID,
		},
	}, nil
}

type PageField2Input struct {
	Name   string
	Update struct {
		Name  *string
		Value *string
	}
}
type updatePageArgs struct {
	ID    gql.ID
	Input *struct {
		Name     *string
		Slug     *string
		ParentID *string
		Fields   *[]PageField2Input
	}
	AddFields    *[]core.PageField
	RemoveFields *[]string
}

func (res *Resolver) UpdatePage(args updatePageArgs) (*pageCreateResultResolver, error) {
	localID, err := fromGlobalID("page", string(args.ID))
	if err != nil {
		return nil, err
	}
	if args.Input != nil || args.AddFields != nil || args.RemoveFields != nil {
		page, err := res.dataSource.GetPage(localID)
		if err != nil {
			return nil, err
		}
		pageInput := core.PageInput{}
		if args.Input != nil {
			if args.Input.Name != nil {
				if *args.Input.Name == "" {
					return &pageCreateResultResolver{
						dataSource: res.dataSource,
						data: pageCreateResult{
							userErrors: &[]userError{
								userError{
									field:   "name",
									message: errNoEmpty("name").Error(),
								},
							},
						},
					}, nil
				}
				pageInput.Name = args.Input.Name
			}
			if args.Input.Slug != nil {
				if *args.Input.Slug == "" {
					return &pageCreateResultResolver{
						dataSource: res.dataSource,
						data: pageCreateResult{
							userErrors: &[]userError{
								userError{
									field:   "slug",
									message: errNoEmpty("slug").Error(),
								},
							},
						},
					}, nil
				}
				pageInput.Slug = args.Input.Slug
			}
			if args.Input.ParentID != nil {
				parentID := args.Input.ParentID
				parentObjectID := bson.ObjectId(*parentID)
				pageInput.ParentID = &parentObjectID
			}
			if args.Input.Fields != nil {
				fields := *args.Input.Fields
				for inputFieldIndex := range fields {
					found := false
					for pageFieldIndex := range page.Fields {
						if page.Fields[pageFieldIndex].Name == fields[inputFieldIndex].Name {
							found = true
							if fields[inputFieldIndex].Update.Name != nil {
								page.Fields[pageFieldIndex].Name = *fields[inputFieldIndex].Update.Name
							}
							if fields[inputFieldIndex].Update.Value != nil {
								page.Fields[pageFieldIndex].Value = *fields[inputFieldIndex].Update.Value
							}
						}
					}
					if !found {
						return &pageCreateResultResolver{
							dataSource: res.dataSource,
							data: pageCreateResult{
								userErrors: &[]userError{
									userError{
										field:   fields[inputFieldIndex].Name,
										message: errNoEmpty(fields[inputFieldIndex].Name).Error(),
									},
								},
							},
						}, nil
					}
				}
				pageInput.Fields = &page.Fields
			}
		}
		if args.AddFields != nil {
			if pageInput.Fields == nil {
				fields := make([]core.PageField, len(page.Fields))
				copy(fields, page.Fields)
				pageInput.Fields = &fields
			}
			tmp := append(*pageInput.Fields, *args.AddFields...)
			pageInput.Fields = &tmp
		}
		if args.RemoveFields != nil {
			if pageInput.Fields == nil {
				fields := make([]core.PageField, len(page.Fields))
				copy(fields, page.Fields)
				pageInput.Fields = &fields
			}
			fields := *pageInput.Fields
			removeFields := *args.RemoveFields
			for inputFieldIndex := range removeFields {
				found := false
				for pageFieldIndex := range page.Fields {
					if fields[pageFieldIndex].Name == removeFields[inputFieldIndex] {
						found = true
					}
				}
				if !found {
					return &pageCreateResultResolver{
						dataSource: res.dataSource,
						data: pageCreateResult{
							userErrors: &[]userError{
								userError{
									field:   fields[inputFieldIndex].Name,
									message: errNoEmpty(fields[inputFieldIndex].Name).Error(),
								},
							},
						},
					}, nil
				}
			}
			fields = []core.PageField{}
			pageFields := *pageInput.Fields
			for pageFieldIndex := range pageFields {
				found := false
				for inputFieldIndex := range removeFields {
					if pageFields[pageFieldIndex].Name == removeFields[inputFieldIndex] {
						found = true
					}
				}
				if !found {
					fields = append(fields, pageFields[pageFieldIndex])
				}
			}
			pageInput.Fields = &fields
		}
		err = res.dataSource.UpdatePage(localID, pageInput)
		if err != nil {
			return nil, err
		}
	}
	page, err := res.dataSource.GetPage(localID)
	if err != nil {
		return nil, err
	}
	return &pageCreateResultResolver{
		dataSource: res.dataSource,
		data: pageCreateResult{
			page: page,
		},
	}, nil
}

type renamePageFieldArgs struct {
	ID    gql.ID
	Input struct {
		Name     string
		RenameTo string
	}
}

func (res *Resolver) RenamePageField(args renamePageFieldArgs) (*pageFieldOperationResultResolver, error) {
	localID, err := fromGlobalID("page", string(args.ID))
	if err != nil {
		return nil, err
	}
	if len(args.Input.RenameTo) == 0 {
		return &pageFieldOperationResultResolver{
			dataSource: res.dataSource,
			data: pageFieldOperationResult{
				userErrors: &[]userError{
					userError{
						field:   "renameTo",
						message: errNoEmpty("renameTo").Error(),
					},
				},
				pageID: localID,
			},
		}, nil
	}
	page, err := res.dataSource.GetPage(localID)
	if err != nil {
		return nil, err
	}
	found := false
	var field core.PageField
	for fieldIndex := range page.Fields {
		if page.Fields[fieldIndex].Name == args.Input.Name {
			found = true
			field = page.Fields[fieldIndex]
		}
	}
	if !found {
		return nil, core.ErrNoField(args.Input.Name)
	}
	field.Name = args.Input.RenameTo
	err = res.dataSource.AddPageField(localID, field)
	if err != nil {
		return nil, err
	}
	err = res.dataSource.RemovePageField(localID, args.Input.Name)
	if err != nil {
		return nil, err
	}
	return &pageFieldOperationResultResolver{
		dataSource: res.dataSource,
		data: pageFieldOperationResult{
			pageID: localID,
		},
	}, nil
}

type updatePageFieldArgs struct {
	ID    gql.ID
	Input struct {
		Name  string
		Value string
	}
}

func (res *Resolver) UpdatePageField(args updatePageFieldArgs) (*pageFieldOperationResultResolver, error) {
	localID, err := fromGlobalID("page", string(args.ID))
	if err != nil {
		return nil, err
	}
	err = res.dataSource.UpdatePageField(localID, args.Input.Name, args.Input.Value)
	if err != nil {
		return nil, err
	}
	return &pageFieldOperationResultResolver{
		dataSource: res.dataSource,
		data: pageFieldOperationResult{
			pageID: localID,
		},
	}, nil
}

type removePageFieldArgs struct {
	ID    gql.ID
	Input struct {
		Name string
	}
}

func (res *Resolver) RemovePageField(args removePageFieldArgs) (*pageFieldOperationResultResolver, error) {
	localID, err := fromGlobalID("page", string(args.ID))
	if err != nil {
		return nil, err
	}
	err = res.dataSource.RemovePageField(localID, args.Input.Name)
	if err != nil {
		return nil, err
	}
	return &pageFieldOperationResultResolver{
		dataSource: res.dataSource,
		data: pageFieldOperationResult{
			pageID: localID,
		},
	}, nil
}

type pageArgs struct {
	ID gql.ID
}

func (res *Resolver) Page(args pageArgs) (*pageResolver, error) {
	localID, err := fromGlobalID("page", string(args.ID))
	if err != nil {
		return nil, err
	}
	result, err := res.dataSource.GetPage(localID)
	if err != nil {
		return nil, err
	}
	return &pageResolver{
		dataSource: res.dataSource,
		data:       &result,
	}, nil
}

type removePageArgs struct {
	ID gql.ID
}

func (res *Resolver) RemovePage(args removePageArgs) (*pageRemoveResultResolver, error) {
	localID, err := fromGlobalID("page", string(args.ID))
	if err != nil {
		return nil, err
	}
	err = res.dataSource.RemovePage(localID)
	if err != nil {
		return nil, err
	}
	return &pageRemoveResultResolver{
		dataSource: res.dataSource,
		data: pageRemoveResult{
			removedObjectID: args.ID,
		},
	}, nil
}
