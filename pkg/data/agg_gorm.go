package data

import (
	"gorm.io/gorm"
)

func dispatchEventsGorm(db *gorm.DB) {
	if agg, ok := isModel[Agg](db); ok {
		err := dispatchEvents(db.Statement.Context, agg)
		if err != nil {
			db.AddError(err)
		}
	}
}

// RegisterAggCallbacks register callback into GORM DB
func RegisterAggCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get("agg:create_events") == nil {
		callback.Create().After("gorm:after_create").Register("agg:create_events", dispatchEventsGorm)
	}
	if callback.Update().Get("agg:update_events") == nil {
		callback.Update().After("gorm:after_update").Register("agg:update_events", dispatchEventsGorm)
	}
	if callback.Delete().Get("agg:delete_events") == nil {
		callback.Update().After("gorm:after_delete").Register("agg:delete_events", dispatchEventsGorm)
	}
}
