package session

import (
	"encoding/json"
	"errors"

	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/types"
	"github.com/s3nt3/sqlvine/internal/logger"
)

type INT64 int64

func (u *INT64) UnmarshalJSON(bs []byte) error {
	var i uint64
	if err := json.Unmarshal(bs, &i); err == nil {
		*u = INT64(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(bs, &s); err != nil {
		return errors.New("expected a string or an integer")
	}
	if err := json.Unmarshal([]byte(s), &i); err != nil {
		return err
	}
	*u = INT64(i)
	return nil
}

type Column struct {
	ID         INT64  `json:"id"`
	Name       string `json:"name"`
	Table      string `json:"table"`
	Type       string `json:"type"`
	PrimaryKey bool   `json:"primary_key"`
}

type Index struct {
	ID      INT64     `json:"id"`
	Name    string    `json:"name"`
	Table   string    `json:"table"`
	Columns []*Column `json:"columns"`
}

type Table struct {
	ID      INT64  `json:"id"`
	Name    string `json:"name"`
	Charset string `json:"charset"`
	Collate string `json:"collate"`

	Columns []*Column `json:"columns"`
	Indices []*Index  `json:"indices"`
}

type Schema struct {
	TableVec []*Table
	TableMap map[string]*Table
}

func NewSchema(s string) *Schema {
	var tables []*Table
	err := json.Unmarshal([]byte(s), &tables)
	if err != nil {
		logger.L.Panic(err.Error())
	}

	schema := &Schema{
		TableVec: []*Table{},
		TableMap: make(map[string]*Table),
	}

	for idx, table := range tables {
		table.ID = INT64(idx)
		schema.TableVec = append(schema.TableVec, table)
		schema.TableMap[table.Name] = table
	}

	return schema
}

func newMySQLTypeLong() types.FieldType {
	return *(types.NewFieldType(mysql.TypeLong))
}

func newMySQLTypeDouble() types.FieldType {
	return *(types.NewFieldType(mysql.TypeDouble))
}

func newMySQLTypeVarchar() (t types.FieldType) {
	defer func() {
		t.Charset, t.Collate = types.DefaultCharsetForType(mysql.TypeVarchar)
	}()

	return *(types.NewFieldType(mysql.TypeVarchar))
}

func newMySQLTypeDate() types.FieldType {
	return *(types.NewFieldType(mysql.TypeDate))
}

func (s *Schema) GetSchemaInfo() (tables []*model.TableInfo) {
	for tidx, table := range s.TableVec {
		tables = append(tables, &model.TableInfo{
			Name: model.NewCIStr(table.Name),
		})

		for cidx, column := range table.Columns {
			tables[tidx].Columns = append(tables[tidx].Columns, &model.ColumnInfo{
				ID:     int64(column.ID),
				Offset: cidx,
				Name:   model.NewCIStr(column.Name),
				State:  model.StatePublic,
			})

			switch column.Type {
			case "float":
				tables[tidx].Columns[cidx].FieldType = newMySQLTypeDouble()
			case "int":
				tables[tidx].Columns[cidx].FieldType = newMySQLTypeLong()
			case "text", "varchar":
				tables[tidx].Columns[cidx].FieldType = newMySQLTypeVarchar()
			case "datetime":
				tables[tidx].Columns[cidx].FieldType = newMySQLTypeDate()
			}

			if column.PrimaryKey {
				tables[tidx].Columns[cidx].Flag = mysql.PriKeyFlag | mysql.NotNullFlag
				tables[tidx].PKIsHandle = (column.Type == "int")
			} else {
				tables[tidx].Columns[cidx].Flag = mysql.NotNullFlag
			}
		}
	}

	return tables
}
