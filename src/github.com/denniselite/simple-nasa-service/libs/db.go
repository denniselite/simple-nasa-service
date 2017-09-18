package libs

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"log"
	"os"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"regexp"
	"time"
	"reflect"
	"github.com/jinzhu/gorm"
	"unicode"
	_ "github.com/lib/pq"
)

type DB struct {
	gorm.DB
}

// Open connection to database
func ConnectDB(dialect string, connectionString string) (*DB, error) {
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		return nil, err
	}
	db.Callback().Update().Replace("gorm:update", ReplaceUpdate)
	return &DB{*db}, nil
}

func ReplaceUpdate(scope *gorm.Scope) {
	if !scope.HasError() {
		var sqls []string

		if updateAttrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
			attrsMap := updateAttrs.(map[string]interface{})
			var keys []string
			for k := range attrsMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, column := range keys {
				sqls = append(sqls, fmt.Sprintf("%v = %v", scope.Quote(column), scope.AddToVars(attrsMap[column])))
			}
		} else {
			for _, field := range scope.Fields() {
				if changeableField(field, scope) {
					if !field.IsPrimaryKey && field.IsNormal {
						sqls = append(sqls, fmt.Sprintf("%v = %v", scope.Quote(field.DBName), scope.AddToVars(field.Field.Interface())))
					} else if relationship := field.Relationship; relationship != nil && relationship.Kind == "belongs_to" {
						for _, foreignKey := range relationship.ForeignDBNames {
							if foreignField, ok := scope.FieldByName(foreignKey); ok && !changeableField(foreignField, scope) {
								sqls = append(sqls,
									fmt.Sprintf("%v = %v", scope.Quote(foreignField.DBName), scope.AddToVars(foreignField.Field.Interface())))
							}
						}
					}
				}
			}
		}

		var extraOption string
		if str, ok := scope.Get("gorm:update_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		if len(sqls) > 0 {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v%v%v",
				scope.QuotedTableName(),
				strings.Join(sqls, ", "),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func changeableField(field *gorm.Field, scope *gorm.Scope) bool {
	if selectAttrs := scope.SelectAttrs(); len(selectAttrs) > 0 {
		for _, attr := range selectAttrs {
			if field.Name == attr || field.DBName == attr {
				return true
			}
		}
		return false
	}

	for _, attr := range scope.OmitAttrs() {
		if field.Name == attr || field.DBName == attr {
			return false
		}
	}

	return true
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

// Debug database with token and own logger
func (db *DB) DebugT(token string) *DB {
	clone := db.New()
	clone.SetLogger(Logger{log.New(os.Stdout, log.Prefix(), log.Flags()), token})
	return &DB{*clone.Debug()}
}

type JsonField map[string]interface{}

func (f JsonField) Value() (driver.Value, error) {
	return JsonValue(f)
}

func (f *JsonField) Scan(src interface{}) error {
	return JsonScan(src, f)
}

func JsonValue(j interface{}) (val driver.Value, err error) {
	val, err = json.Marshal(j)
	return
}

func JsonScan(src interface{}, out interface{}) (err error) {
	source, ok := src.([]byte)
	if !ok {
		err = errors.New("Incorrect field type. []byte expected.")
		return
	}

	err = json.Unmarshal(source, &out)
	return
}

var (
	sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)
)

// LogWriter log writer interface
type LogWriter interface {
	Println(v ...interface{})
	Output(calldepth int, s string) error
}

type Logger struct {
	LogWriter
	Token string
}

// Print format & print log
func (logger Logger) Print(values ...interface{}) {
	if len(values) > 1 {
		level := values[0]
		source := values[1]

		messages := []interface{}{logger.Token}

		if level == "sql" {
			// duration
			messages = append(messages, fmt.Sprintf(" [%.2fms] ", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0))
			// sql
			var sql string
			var formattedValues []string

			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}

			var formattedValuesLength = len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}

			messages = append(messages, sql)
		} else {
			messages = append(messages, "\033[31;1m")
			messages = append(messages, values[2:]...)
			messages = append(messages, "\033[0m")
		}

		// Определяем вложенность для вывода имени файла в лог
		calldepth := 2
		for ; calldepth < 15; calldepth++ {
			_, file, line, ok := runtime.Caller(calldepth)
			if ok && fmt.Sprintf("%v:%v", file, line) == source {
				break
			}
		}

		logger.Output(calldepth+1, fmt.Sprint(messages...))
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}