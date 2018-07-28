package mock

import (
	"fmt"

	"github.com/dominik-zeglen/inkster/core"
	"github.com/globalsign/mgo/bson"
	"github.com/gosimple/slug"
)

func (adapter Adapter) findPage(id *bson.ObjectId, slug *string) (int, error) {
	if id != nil {
		for index := range pages {
			if pages[index].ID == *id {
				return index, nil
			}
		}
		return 0, fmt.Errorf("Page %s does not exist", id)
	}
	if slug != nil {
		for index := range pages {
			if pages[index].Slug == *slug {
				return index, nil
			}
		}
		return 0, fmt.Errorf("Page %s does not exist", *slug)
	}
	if id == nil && slug == nil {
		return 0, fmt.Errorf("findPage() requires at least one argument")
	}
	return 0, fmt.Errorf("")
}
func (adapter Adapter) findPageField(id bson.ObjectId, name string) (int, int, error) {
	index, err := adapter.findPage(&id, nil)
	if err != nil {
		return 0, 0, err
	}
	for fieldIndex := range pages[index].Fields {
		if pages[index].Fields[fieldIndex].Name == name {
			return index, fieldIndex, nil
		}
	}
	return 0, 0, core.ErrNoField(name)
}

// AddPage puts page in the database
func (adapter Adapter) AddPage(page core.Page) (core.Page, error) {
	err := page.Validate()
	if err != nil {
		return core.Page{}, err
	}
	_, err = adapter.findTemplate(nil, &page.Name)
	if err == nil {
		return core.Page{}, core.ErrPageExists(page.Name)
	}
	if page.ID == "" {
		page.ID = bson.NewObjectId()
	} else {
		_, err = adapter.findPage(&page.ID, nil)
		if err == nil {
			return core.Page{}, core.ErrPageExists(page.ID.String())
		}
	}
	if page.Slug == "" {
		slug := slug.Make(page.Name)
		page.Slug = slug
	}
	pages = append(pages, page)
	return page, nil
}

// AddPageFromTemplate creates new page based on a chosen template
func (adapter Adapter) AddPageFromTemplate(
	page core.PageInput,
	templateID bson.ObjectId,
) (core.Page, error) {
	template, err := adapter.GetTemplate(templateID)
	if err != nil {
		return core.Page{}, err
	}
	var fields []core.PageField
	for _, field := range template.Fields {
		fields = append(fields, core.PageField{
			Name:  field.Name,
			Type:  field.Type,
			Value: "",
		})
	}
	if page.Name == nil {
		return core.Page{}, core.ErrNoEmpty("name")
	}
	if page.ParentID == nil {
		return core.Page{}, core.ErrNoEmpty("parentID")
	}
	inputPage := core.Page{
		Name:     *page.Name,
		ParentID: *page.ParentID,
		Fields:   fields,
	}
	if page.Slug != nil {
		inputPage.Slug = *page.Slug
	} else {
		slug := slug.Make(*page.Name)
		inputPage.Slug = slug
	}
	return adapter.AddPage(inputPage)
}

// AddPageField adds to page a new field at the end of it's field list
func (adapter Adapter) AddPageField(pageID bson.ObjectId, field core.PageField) error {
	err := field.Validate()
	if err != nil {
		return err
	}

	index, _, err := adapter.findPageField(pageID, field.Name)
	if err == nil {
		return core.ErrFieldExists(field.Name)
	}
	pages[index].Fields = append(pages[index].Fields, field)
	return nil
}

// GetPage allows user to fetch page by ID from database
func (adapter Adapter) GetPage(id bson.ObjectId) (core.Page, error) {
	index, err := adapter.findPage(&id, nil)
	return pages[index], err
}

// GetPageBySlug allows user to fetch page by slug from database
func (adapter Adapter) GetPageBySlug(slug string) (core.Page, error) {
	index, err := adapter.findPage(nil, &slug)
	return pages[index], err
}

// GetPagesFromDirectory allows user to fetch pages by their parentId from database
func (adapter Adapter) GetPagesFromDirectory(id bson.ObjectId) ([]core.Page, error) {
	var returnPages []core.Page
	for index := range pages {
		if pages[index].ParentID == id {
			returnPages = append(returnPages, pages[index])
		}
	}
	return returnPages, nil
}

// UpdatePage allows user to update page properties
func (adapter Adapter) UpdatePage(pageID bson.ObjectId, data core.PageInput) error {
	index, err := adapter.findPage(&pageID, nil)
	if err != nil {
		return err
	}
	if data.Name != nil {
		pages[index].Name = *data.Name
	}
	if data.Slug != nil {
		i, err := adapter.findPage(nil, data.Slug)
		if i != index && err == nil {
			return core.ErrPageExists(*data.Slug)
		}
		pages[index].Slug = *data.Slug
	}
	if data.ParentID != nil {
		_, err = adapter.findDirectory(*data.ParentID)
		if err == nil {
			return err
		}
		pages[index].ParentID = *data.ParentID
	}
	if data.Fields != nil {
		fields := make([]core.PageField, len(*data.Fields))
		copy(fields, *data.Fields)
		pages[index].Fields = fields
	}
	return nil
}

// UpdatePageField removes field from page
func (adapter Adapter) UpdatePageField(pageID bson.ObjectId, pageFieldName string, data string) error {
	index, fieldIndex, err := adapter.findPageField(pageID, pageFieldName)
	if err != nil {
		return err
	}
	fields := make([]core.PageField, len(pages[index].Fields))
	copy(fields, pages[index].Fields)
	fields[fieldIndex].Value = data
	pages[index].Fields = fields
	return nil
}

// RemovePage removes page from database
func (adapter Adapter) RemovePage(pageID bson.ObjectId) error {
	index, err := adapter.findPage(&pageID, nil)
	if err != nil {
		return err
	}
	pages = append(pages[:index], pages[index+1:]...)
	return nil
}

// RemovePageField removes field from page
func (adapter Adapter) RemovePageField(pageID bson.ObjectId, pageFieldName string) error {
	index, fieldIndex, err := adapter.findPageField(pageID, pageFieldName)
	if err != nil {
		return err
	}
	fields := make([]core.PageField, len(pages[index].Fields))
	copy(fields, pages[index].Fields)
	pages[index].Fields = append(
		fields[:fieldIndex],
		fields[fieldIndex+1:]...,
	)
	return nil
}
