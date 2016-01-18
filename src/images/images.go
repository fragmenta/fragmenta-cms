// Package images represents image resources
package images

import (
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/fragmenta/model"
	"github.com/fragmenta/model/file"
	"github.com/fragmenta/model/validate"
	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// Image represents an image on disk
type Image struct {
	model.Model
	status.ModelStatus
	Name string
	Path string
	Sort int64
}

// AllowedParams returns the params this model allows
func AllowedParams() []string {
	return []string{"status", "name", "path", "sort"}
}

// NewWithColumns creates a image from database columns - used by query in creating objects
func NewWithColumns(cols map[string]interface{}) *Image {

	image := New()
	image.Id = validate.Int(cols["id"])
	image.CreatedAt = validate.Time(cols["created_at"])
	image.UpdatedAt = validate.Time(cols["updated_at"])
	image.Status = validate.Int(cols["status"])
	image.Name = validate.String(cols["name"])
	image.Path = validate.String(cols["path"])
	image.Sort = validate.Int(cols["sort"])

	return image
}

// New sets up a new image with default values
func New() *Image {
	image := &Image{}
	image.Model.Init()
	image.TableName = "images"
	image.Status = status.Published
	image.Sort = 1

	return image
}

// SaveImageRepresentations saves files to disk
func (m *Image) SaveImageRepresentations(f multipart.File) error {

	// If we have no path, set it to a default value /files/images/id/name
	if len(m.Path) == 0 {
		err := m.SetDefaultOriginalPath()
		if err != nil {
			return err
		}
	}

	// Write out several representations of this file to disk
	options := []file.Options{
		file.Options{Path: path.Join("public", m.Path), MaxHeight: 4000, MaxWidth: 4000, Quality: 100},
		file.Options{Path: path.Join("public", m.LargePath()), MaxHeight: 2000, MaxWidth: 2000, Quality: 70},
		file.Options{Path: path.Join("public", m.SmallPath()), MaxHeight: 400, MaxWidth: 400, Quality: 60},
		file.Options{Path: path.Join("public", m.IconPath()), MaxHeight: 200, MaxWidth: 200, Quality: 60},
	}

	// Make sure our path exists first
	err := file.CreatePathTo(path.Join("public", m.Path))
	if err != nil {
		return err
	}

	return file.SaveJpegRepresentations(f, options)
}

// SetDefaultOriginalPath assumes that the path ends in .jpg.
func (m *Image) SetDefaultOriginalPath() error {
	m.Path = fmt.Sprintf("files/images/%d/%s.jpg", m.Id, file.SanitizeName(m.Name))
	return m.Update(map[string]string{"path": m.Path})
}

// LargePath returns the image large path.
func (m *Image) LargePath() string {
	return strings.Replace(m.Path, ".jpg", "-large.jpg", -1)
}

// SmallPath returns the image small path.
func (m *Image) SmallPath() string {
	return strings.Replace(m.Path, ".jpg", "-small.jpg", -1)
}

// IconPath returns the image icon path.
func (m *Image) IconPath() string {
	return strings.Replace(m.Path, ".jpg", "-icon.jpg", -1)
}

// Update this image.
func (m *Image) Update(params map[string]string) error {

	// Remove params not in AllowedParams.
	params = model.CleanParams(params, AllowedParams())

	err := validateParams(params)
	if err != nil {
		return err
	}

	// Make sure updated_at is set to the current time
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("id=?", m.Id).Update(params)
}

// Destroy this image.
func (m *Image) Destroy() error {
	// We also need to delete the image files on disk
	// TODO - delete image files on disk after image destroy.

	return Query().Where("id=?", m.Id).Delete()
}

// OwnedBy images permissions require a join table with users, updated on create.
func (m *Image) OwnedBy(id int64) bool {
	return true
}

// validateParams checks these parameters conform to expectations.
func validateParams(unsafeParams map[string]string) error {

	// Now check params are as we expect,
	err := validate.Length(unsafeParams["id"], 0, -1)
	if err != nil {
		return err
	}

	return err
}

// Create inserts a new image record in the database and returns the ID.
func Create(params map[string]string, fh *multipart.FileHeader) (int64, error) {

	// Remove params not in AllowedParams
	params = model.CleanParams(params, AllowedParams())

	err := validateParams(params)
	if err != nil {
		return 0, err
	}

	// Update/add some params by default.
	params["created_at"] = query.TimeString(time.Now().UTC())
	params["updated_at"] = query.TimeString(time.Now().UTC())

	id, err := Query().Insert(params)

	if fh != nil && id != 0 {
		// Retreive the form image data by opening the referenced tmp file.
		f, err := fh.Open()
		if err != nil {
			return id, err
		}

		// Now retrieve the image concerned, and save the file representations.
		image, err := Find(id)
		if err != nil {
			return id, err
		}

		// Save files to disk using the passed in file data (if any).
		err = image.SaveImageRepresentations(f)
		if err != nil {
			return id, err
		}
	}

	return id, err
}

// Query creates a new query relation referencing this model.
func Query() *query.Query {
	return query.New("images", "id")
}

// Ordered returns an ordered result set.
func Ordered() *query.Query {
	return Query().Order("images.sort asc")
}

// Where returns a query shortcut for the common where query on images.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Find fetches a single record by id in params.
func Find(id int64) (*Image, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all results for this query.
func FindAll(q *query.Query) ([]*Image, error) {

	// Fetch query.Results from query.
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of pages constructed from the results.
	var imageList []*Image
	for _, r := range results {
		image := NewWithColumns(r)
		imageList = append(imageList, image)
	}

	return imageList, nil
}
