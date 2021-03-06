package schemamanagement

import (
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/timescale/outflux/internal/connections"
	influxSchema "github.com/timescale/outflux/internal/schemamanagement/influx"
	"github.com/timescale/outflux/internal/schemamanagement/influx/discovery"
	tsSchema "github.com/timescale/outflux/internal/schemamanagement/ts"
)

// SchemaManagerService defines methods for creating SchemaManagers
type SchemaManagerService interface {
	Influx(client influx.Client, db, rp string, onConflictConvertIntToFloat bool) SchemaManager
	TimeScale(dbConn connections.PgxWrap, schema, chunkTimeInterval string) SchemaManager
}

// NewSchemaManagerService returns an instance of SchemaManagerService
func NewSchemaManagerService(measureExplorer discovery.MeasureExplorer, tagExplorer discovery.TagExplorer, fieldExplorer discovery.FieldExplorer) SchemaManagerService {
	return &schemaManagerService{
		tagExplorer:     tagExplorer,
		fieldExplorer:   fieldExplorer,
		measureExplorer: measureExplorer,
	}
}

type schemaManagerService struct {
	tagExplorer     discovery.TagExplorer
	fieldExplorer   discovery.FieldExplorer
	measureExplorer discovery.MeasureExplorer
}

// Influx creates new schema manager that can discover influx data sets
func (s *schemaManagerService) Influx(client influx.Client, db, rp string, onConflictConvertIntToFloat bool) SchemaManager {
	return influxSchema.NewSchemaManager(client, db, rp, onConflictConvertIntToFloat, s.measureExplorer, s.tagExplorer, s.fieldExplorer)
}

func (s *schemaManagerService) TimeScale(dbConn connections.PgxWrap, schema, chunkTimeInterval string) SchemaManager {
	return tsSchema.NewTSSchemaManager(dbConn, schema, chunkTimeInterval)
}
