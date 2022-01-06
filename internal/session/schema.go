package session

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/types"
	"github.com/s3nt3/sqlvine/internal/logger"
	"github.com/s3nt3/sqlvine/internal/util"
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
	ID          INT64  `json:"id"`
	Name        string `json:"name"`
	OriginName  string
	Table       string `json:"table"`
	OriginTable string
	Type        string `json:"type"`
	Size        INT64  `json:"size"`
	PrimaryKey  bool   `json:"primary_key"`
}

type Index struct {
	ID          INT64  `json:"id"`
	Name        string `json:"name"`
	OriginName  string
	Table       string `json:"table"`
	OriginTable string
	Columns     []*Column `json:"columns"`
}

type Table struct {
	ID         INT64  `json:"id"`
	Name       string `json:"name"`
	OriginName string
	Charset    string `json:"charset"`
	Collate    string `json:"collate"`

	Columns []*Column `json:"columns"`
	Indices []*Index  `json:"indices"`

	*util.Random
}

func NewTable() *Table {
	return &Table{
		Charset: "utf8mb4",
		Collate: "utf8mb4_bin",
		Columns: []*Column{},
		Indices: []*Index{},
		Random:  util.NewRandom(),
	}
}

func MergeTable(l *Table, r *Table) *Table {
	if l != nil {
		if r != nil {
			t := NewTable()
			for _, column := range l.Columns {
				column.Name = fmt.Sprintf("c%d", len(t.Columns))
				t.Columns = append(t.Columns, column)
			}

			for _, column := range r.Columns {
				column.Name = fmt.Sprintf("c%d", len(t.Columns))
				t.Columns = append(t.Columns, column)
			}

			for _, indice := range l.Indices {
				indice.Name = fmt.Sprintf("c%d", len(t.Columns))
				t.Indices = append(t.Indices, indice)
			}

			for _, indice := range r.Indices {
				indice.Name = fmt.Sprintf("c%d", len(t.Columns))
				t.Indices = append(t.Indices, indice)
			}

			return t
		}

		return l
	}

	return r
}

func (t *Table) GetRandomColumn() *Column {
	if len(t.Columns) > 0 {
		return t.Columns[t.RandomNum(len(t.Columns))]
	}

	return nil
}

func (t *Table) SetID(id int64) {
	t.ID = INT64(id)
}

func (t *Table) ReName(name string) {
	t.OriginName = t.Name
	t.Name = name

	for _, column := range t.Columns {
		column.OriginTable = column.Table
		column.Table = name
	}
}

type Schema struct {
	*util.Random

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
		Random:   util.NewRandom(),
		TableVec: []*Table{},
		TableMap: make(map[string]*Table),
	}

	for idx, table := range tables {
		table.ID = INT64(idx)
		table.Random = util.NewRandom()

		schema.TableVec = append(schema.TableVec, table)
		schema.TableMap[table.Name] = table

		for _, columns := range table.Columns {
			columns.Table = table.Name
		}

		for _, indice := range table.Indices {
			indice.Table = table.Name
		}
	}

	return schema
}

func (s *Schema) AddTable(t *Table) {
	if t != nil {
		s.TableVec = append(s.TableVec, t)
		s.TableMap[t.Name] = t
	}
}

func (s *Schema) GetRandomTable() *Table {
	if len(s.TableVec) > 0 {
		return s.TableVec[s.RandomNum(len(s.TableVec))]
	}

	return nil
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
