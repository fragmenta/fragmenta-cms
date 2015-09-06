// The [[ .fragmenta_resources ]] package
package [[ .fragmenta_resources ]]

import (
    "time"
    
	"github.com/fragmenta/query"
	"github.com/fragmenta/model"
    "github.com/fragmenta/model/validate"
    "github.com/fragmenta/model/status"
)

// The [[ .fragmenta_resources ]] model type
type [[ .Fragmenta_Resource ]] struct {
	model.Model
    status.ModelStatus
[[ .fragmenta_fields ]]
}

// Create a [[ .fragmenta_resource ]] from database columns - used by query in creating objects
func (m *[[ .Fragmenta_Resource ]]) New(cols map[string]interface{}) *[[ .Fragmenta_Resource ]] {

	[[ .fragmenta_resource ]] := New()
	[[ .fragmenta_resource ]].Id = validate.Int(cols["id"])
	[[ .fragmenta_resource ]].CreatedAt = validate.Time(cols["created_at"])
	[[ .fragmenta_resource ]].UpdatedAt = validate.Time(cols["updated_at"])
    [[ .fragmenta_resource ]].Status = validate.Int(cols["status"])
[[ .fragmenta_new_fields ]]
    
	return [[ .fragmenta_resource ]]
}

// Update this [[ .fragmenta_resource ]]
func (m *[[ .Fragmenta_Resource ]]) Update(params map[string]string) error {

	err := ValidateParams(params)
	if err != nil {
		return err
	}

    // Make sure updated_at is set to the current time
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("id=?", m.Id).Update(params)
}

// Delete this [[ .fragmenta_resource ]]
func (m *[[ .Fragmenta_Resource ]]) Destroy() error {
	return Query().Where("id=?", m.Id).Delete()
}

// Which parameters does this model allow to be edited in forms?
func AcceptedParams() []string {
	return []string{"id","status",[[ .fragmenta_columns ]]}
}

// Check these parameters conform to AcceptedParams, and pass validation
func ValidateParams(unsafeParams map[string]string) error {

	// First check for params we don't accept - we fail if we receive columns we don't expect
	_, err := validate.CleanParams(unsafeParams, AcceptedParams())
	if err != nil {
		return err
	}

	// Now check params are as we expect
	err = validate.Length(unsafeParams["id"], 0, -1)
	if err != nil {
		return err
	}

	return err
}

// Set up a new [[ .fragmenta_resource ]] with default values
func New() *[[ .Fragmenta_Resource ]] {
	[[ .fragmenta_resource ]] := &[[ .Fragmenta_Resource ]]{}
	[[ .fragmenta_resource ]].Model.Init()
	[[ .fragmenta_resource ]].TableName = "[[ .fragmenta_resources ]]"
    [[ .fragmenta_resource ]].Status = status.Draft
	
	return [[ .fragmenta_resource ]]
}

// Insert a new [[ .fragmenta_resource ]]
func Create(params map[string]string) (int64, error) {
	err := ValidateParams(params)
	if err != nil {
		return 0, err
	}

	// Update/add some params by default
	params["created_at"] = query.TimeString(time.Now().UTC())
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

// Create a new query relation referencing this model
func Query() *query.Query {
	return query.New("[[ .fragmenta_resources ]]","id")
}

// A shortcut for the common where query on [[ .fragmenta_resources ]]
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Request a single record by id in params
func Find(id int64) (*[[ .Fragmenta_Resource ]], error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return New().New(result), nil
}


