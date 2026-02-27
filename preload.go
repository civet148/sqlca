package sqlca

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/civet148/sqlca/v3/types"
)

func (e *Engine) addPreload(query string, args ...any) *Engine {
	if e.preloads == nil {
		e.preloads = make(map[string][]any)
	}
	e.preloads[query] = args
	return e
}

func (e *Engine) handlePreloads() (err error) {
	for query, args := range e.preloads {
		if err = e.execPreload(query, args...); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) execPreload(query string, args ...any) (err error) {
	if e.model == nil {
		return fmt.Errorf("model is nil, please call Model method first")
	}

	// 获取当前模型类型和值
	modelType := reflect.TypeOf(e.model)
	modelValue := reflect.ValueOf(e.model)

	// 处理指针类型
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}

	// 如果是切片类型，则处理每个元素
	if modelType.Kind() == reflect.Slice {
		for i := 0; i < modelValue.Len(); i++ {
			element := modelValue.Index(i)
			if element.Kind() == reflect.Ptr {
				element = element.Elem()
			}
			err = e.loadAssociations(element, query, args...)
			if err != nil {
				return err
			}
		}
	} else {
		// 单个结构体实例
		err = e.loadAssociations(modelValue, query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// loadAssociations 加载关联数据
func (e *Engine) loadAssociations(modelValue reflect.Value, associationName string, args ...any) error {
	modelType := modelValue.Type()

	// 遍历结构体字段寻找匹配的关联字段
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		// 检查字段名是否匹配关联名称
		if field.Name != associationName {
			continue
		}

		// 获取GORM标签
		gormTag := field.Tag.Get("gorm")
		if gormTag == "" {
			continue
		}

		// 解析GORM标签设置
		settings := parseTagSetting(gormTag, ";")

		// 检查是否为many2many关联
		if junctionTable, ok := settings["MANY2MANY"]; ok {
			return e.loadMany2ManyAssociation(fieldValue, field.Type, junctionTable, args...)
		}

		// 检查是否为foreignKey关联
		if foreignKey, ok := settings["FOREIGNKEY"]; ok {
			return e.loadForeignKeyAssociation(fieldValue, field.Type, foreignKey, args...)
		}

		// 检查是否有references引用
		if references, ok := settings["REFERENCES"]; ok {
			foreignKey, _ := settings["FOREIGNKEY"]
			return e.loadReferenceAssociation(fieldValue, field.Type, foreignKey, references, args...)
		}
	}

	return fmt.Errorf("association '%s' not found", associationName)
}

// loadMany2ManyAssociation 加载多对多关联
func (e *Engine) loadMany2ManyAssociation(fieldValue reflect.Value, fieldType reflect.Type, junctionTable string, args ...any) error {
	if fieldValue.Kind() != reflect.Ptr && fieldValue.Kind() != reflect.Slice {
		return fmt.Errorf("many2many association field must be a pointer to slice")
	}

	// 如果是指针，获取其指向的切片
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			newSlice := reflect.New(fieldType.Elem())
			fieldValue.Set(newSlice)
			fieldValue = newSlice.Elem()
		} else {
			fieldValue = fieldValue.Elem()
		}
	}

	if fieldValue.Kind() != reflect.Slice {
		return fmt.Errorf("many2many association field must be a pointer to slice")
	}

	// 获取主模型的主键值
	mainModelValue := e.model
	mainModelType := reflect.TypeOf(mainModelValue)
	mainModelReflectValue := reflect.ValueOf(mainModelValue)

	// 处理指针
	if mainModelType.Kind() == reflect.Ptr {
		mainModelType = mainModelType.Elem()
		mainModelReflectValue = mainModelReflectValue.Elem()
	}

	// 检查是否是切片类型，如果是则取第一个元素
	if mainModelType.Kind() == reflect.Slice {
		if mainModelReflectValue.Len() == 0 {
			return fmt.Errorf("main model slice is empty, cannot get primary key value")
		}
		// 取第一个元素
		firstElement := mainModelReflectValue.Index(0)
		if firstElement.Kind() == reflect.Ptr {
			if firstElement.IsNil() {
				return fmt.Errorf("first element in main model slice is nil")
			}
			firstElement = firstElement.Elem()
		}
		mainModelType = firstElement.Type()
		mainModelReflectValue = firstElement
	}

	// 尝试找到主键字段 (通常名为ID)
	var mainPKValue interface{}
	for i := 0; i < mainModelType.NumField(); i++ {
		field := mainModelType.Field(i)
		gormTag := field.Tag.Get("gorm")
		settings := parseTagSetting(gormTag, ";")
		if _, isPK := settings["PRIMARYKEY"]; isPK || strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
			fieldVal := mainModelReflectValue.Field(i)
			mainPKValue = fieldVal.Interface()
			break
		}
	}

	if mainPKValue == nil {
		// 如果没找到明确的主键，尝试获取第一个整数类型的字段作为主键
		for i := 0; i < mainModelType.NumField(); i++ {
			field := mainModelType.Field(i)
			fieldVal := mainModelReflectValue.Field(i)
			if strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
				mainPKValue = fieldVal.Interface()
				break
			}
		}
	}

	if mainPKValue == nil {
		return fmt.Errorf("could not find primary key value for main model")
	}

	// 构建关联查询：从junctionTable中查找对应记录，再关联到目标表
	elementType := fieldValue.Type().Elem()
	if elementType.Kind() == reflect.Ptr {
		elementType = elementType.Elem()
	}

	// 创建关联模型实例
	_ = reflect.New(elementType).Interface()

	// 执行查询：SELECT * FROM assoc_table WHERE id IN (SELECT assoc_id FROM junction_table WHERE main_id = ?)
	// 首先从junctionTable获取关联ID列表
	// 解析主模型和关联模型的字段名
	mainModelName := getTableNameFromType(mainModelType)
	assocModelName := getTableNameFromType(elementType)

	// 生成可能的外键列名
	junctionPK1 := mainModelName + "_id"
	junctionPK2 := assocModelName + "_id"

	// 尝试从GORM标签中获取更准确的外键名称
	// 如果主模型是User，关联模型是Role，则junctionPK1=user_id, junctionPK2=role_id
	junctionPK1 = snakeCase(mainModelType.Name()) + "_id"
	junctionPK2 = snakeCase(elementType.Name()) + "_id"

	var assocIDs []interface{}
	junctionQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", junctionPK2, junctionTable, junctionPK1)
	rows, err := e.db.Query(junctionQuery, mainPKValue)
	if err != nil {
		return fmt.Errorf("failed to query junction table: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var assocID interface{}
		if err := rows.Scan(&assocID); err != nil {
			return fmt.Errorf("failed to scan junction row: %v", err)
		}
		assocIDs = append(assocIDs, assocID)
	}

	if len(assocIDs) == 0 {
		// 没有关联的数据，设置空切片
		fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
		return nil
	}

	// 使用IN查询获取关联数据
	// 合并关联ID和额外的查询条件
	totalArgs := make([]interface{}, 0, len(assocIDs)+len(args))
	totalArgs = append(totalArgs, assocIDs...)
	totalArgs = append(totalArgs, args...)

	placeholders := make([]string, len(assocIDs))
	queryArgs := make([]interface{}, len(assocIDs))
	for i, id := range assocIDs {
		placeholders[i] = "?"
		queryArgs[i] = id
	}

	tableName := getTableNameFromType(elementType)

	// 动态构建选择字段列表，而不是使用 *
	columns := getStructFields(elementType, e)
	selectFields := strings.Join(columns, ", ")

	inClause := strings.Join(placeholders, ", ")

	// 构建基础查询
	finalQuery := fmt.Sprintf("SELECT %s FROM %s WHERE id IN (%s)", selectFields, tableName, inClause)

	// 如果有额外的查询条件，则添加它们
	if len(args) > 0 {
		// 支持GORM风格的条件参数，例如: "status = ? AND deleted_at is NULL", 1
		// 或者复合条件: "status = ? AND created_at > ?", "active", "2023-01-01"
		// 检查第一个参数是否为字符串（可能是条件表达式）
		if condition, ok := args[0].(string); ok {
			// 添加条件到查询中
			finalQuery += " AND " + condition
			// 将其余参数添加到查询参数中
			if len(args) > 1 {
				queryArgs = append(queryArgs, args[1:]...)
			}
		}
	}

	// 替换占位符为正确的数据库方言格式
	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		// 重新构建查询以适应PostgreSQL占位符
		// 为IN子句中的ID创建占位符
		inPlace := make([]string, len(assocIDs))
		for i := 0; i < len(assocIDs); i++ {
			inPlace[i] = "$" + strconv.Itoa(i+1)
		}
		inClause = strings.Join(inPlace, ", ")

		// 重新构建基础查询
		finalQuery = fmt.Sprintf("SELECT %s FROM %s WHERE id IN (%s)", selectFields, tableName, inClause)

		// 添加额外条件
		argOffset := len(assocIDs) + 1
		if len(args) > 0 {
			// 检查第一个参数是否为字符串（可能是条件表达式）
			if condition, ok := args[0].(string); ok {
				// 处理包含占位符的字符串条件
				// 计算条件字符串中占位符的数量，并调整argOffset
				placeholderCount := strings.Count(condition, "?")
				finalQuery += " AND " + strings.ReplaceAll(condition, "?", fmt.Sprintf("$%d", argOffset))

				// 更新argOffset为占位符之后的位置
				argOffset += placeholderCount
			} else if len(args)%2 == 0 {
				for i := 0; i < len(args); i += 2 {
					finalQuery += " AND "

					field, ok := args[i].(string)
					if !ok {
						return fmt.Errorf("invalid field name type: %T", args[i])
					}

					operator := "="
					valueIndex := i + 1

					// 检查下一个参数是否是操作符
					if valueIndex < len(args) {
						nextArg, isOp := args[valueIndex].(string)
						if isOp && isInOperator(nextArg) {
							operator = nextArg
							valueIndex++ // 如果下一个是操作符，则实际值是再下一个
						}
					}

					finalQuery += fmt.Sprintf("%s %s $%d", field, operator, argOffset)
					argOffset++
				}
			}
		}

		// 为PostgreSQL重新构建参数列表
		queryArgs = make([]interface{}, 0, len(assocIDs)+len(args))
		queryArgs = append(queryArgs, assocIDs...)
		if len(args) > 0 {
			// 检查第一个参数是否为字符串（可能是条件表达式）
			if _, ok := args[0].(string); ok {
				// 如果是字符串条件表达式，直接将剩余参数添加到查询参数中
				if len(args) > 1 {
					queryArgs = append(queryArgs, args[1:]...)
				}
			} else {
				// 如果不是字符串条件表达式，按照键值对方式处理
				if len(args)%2 == 0 {
					for i := 0; i < len(args); i += 2 {
						// 跳过操作符，只添加值
						valueIndex := i + 1
						// 检查下一个参数是否是操作符
						if valueIndex < len(args) {
							nextArg, isOp := args[valueIndex].(string)
							if isOp && isInOperator(nextArg) {
								valueIndex++ // 如果下一个是操作符，则实际值是再下一个
							}
						}
						if valueIndex < len(args) {
							queryArgs = append(queryArgs, args[valueIndex])
						}
					}
				} else {
					// 如果参数不是偶数个，返回错误
					return fmt.Errorf("unsupported argument format for additional conditions")
				}
			}
		}

	}

	// 查询关联数据
	resultSlice := reflect.New(reflect.SliceOf(reflect.PtrTo(elementType)))

	rows, err = e.db.Query(finalQuery, queryArgs...)
	if err != nil {
		return fmt.Errorf("failed to query associated data: %v", err)
	}
	defer rows.Close()

	// 使用反射将结果填充到切片
	resultSliceVal := resultSlice.Elem()
	for rows.Next() {
		// 获取列信息
		cols, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to get columns: %v", err)
		}

		// 为每一列创建一个接口值来接收数据
		values := make([]interface{}, len(cols))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			return fmt.Errorf("failed to scan associated row: %v", err)
		}

		// 创建新的元素实例
		element := reflect.New(elementType).Elem()

		// 将扫描的结果映射到结构体字段
		colMap := make(map[string]interface{})
		for i, col := range cols {
			colMap[col] = *(values[i].(*interface{}))
		}

		// 使用fetcher将数据映射到结构体
		fetcher := &Fetcher{mapValues: convertMapInterfaceToMapString(colMap)}
		if err := e.fetchToStruct(fetcher, elementType, element); err != nil {
			return fmt.Errorf("failed to map fetched data to struct: %v", err)
		}

		resultSliceVal = reflect.Append(resultSliceVal, element)
	}

	// 设置关联字段值
	fieldValue.Set(resultSliceVal)

	return nil
}

// loadForeignKeyAssociation 加载外键关联
func (e *Engine) loadForeignKeyAssociation(fieldValue reflect.Value, fieldType reflect.Type, foreignKey string, args ...any) error {
	// 获取外键字段的值
	mainModelValue := e.model
	mainModelType := reflect.TypeOf(mainModelValue)
	mainModelReflectValue := reflect.ValueOf(mainModelValue)

	// 处理指针
	if mainModelType.Kind() == reflect.Ptr {
		mainModelType = mainModelType.Elem()
		mainModelReflectValue = mainModelReflectValue.Elem()
	}

	// 检查是否是切片类型，如果是则取第一个元素
	if mainModelType.Kind() == reflect.Slice {
		if mainModelReflectValue.Len() == 0 {
			return fmt.Errorf("main model slice is empty, cannot get foreign key value")
		}
		// 取第一个元素
		firstElement := mainModelReflectValue.Index(0)
		if firstElement.Kind() == reflect.Ptr {
			if firstElement.IsNil() {
				return fmt.Errorf("first element in main model slice is nil")
			}
			firstElement = firstElement.Elem()
		}
		mainModelType = firstElement.Type()
		mainModelReflectValue = firstElement
	}

	// 查找主模型的主键字段值
	var fkValue interface{}
	found := false
	for i := 0; i < mainModelType.NumField(); i++ {
		field := mainModelType.Field(i)
		gormTag := field.Tag.Get("gorm")
		settings := parseTagSetting(gormTag, ";")
		_, isPK := settings["PRIMARYKEY"]
		if isPK || strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
			fieldVal := mainModelReflectValue.Field(i)
			fkValue = fieldVal.Interface()
			found = true
			break
		}
	}

	if !found {
		// 如果没找到明确的主键，尝试获取第一个整数类型的字段作为主键
		for i := 0; i < mainModelType.NumField(); i++ {
			field := mainModelType.Field(i)
			fieldVal := mainModelReflectValue.Field(i)
			if strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
				fkValue = fieldVal.Interface()
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("could not find primary key value for main model")
	}

	// 准备关联模型
	var assocModelValue reflect.Value
	if fieldType.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			assocModelValue = reflect.New(fieldType.Elem()).Elem()
			fieldValue.Set(reflect.New(fieldType.Elem()))
			fieldValue = fieldValue.Elem()
		} else {
			assocModelValue = fieldValue.Elem()
		}
	} else {
		assocModelValue = fieldValue
	}

	// 构建查询：SELECT {fields} FROM assoc_table WHERE id = fk_value
	elementType := fieldType
	if fieldType.Kind() == reflect.Ptr {
		elementType = fieldType.Elem()
	}

	tableName := getTableNameFromType(elementType)

	// 动态构建选择字段列表，而不是使用 *
	columns := getStructFields(elementType, e)
	selectFields := strings.Join(columns, ", ")

	// 将 foreignKey 转换为蛇形命名法，以匹配数据库中的列名
	foreignKeySnake := snakeCase(foreignKey)

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", selectFields, tableName, foreignKeySnake)

	// 处理不同数据库的占位符
	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", selectFields, tableName, foreignKeySnake)
	}

	rows, err := e.db.Query(query, fkValue)
	if err != nil {
		return fmt.Errorf("failed to query foreign key association: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		// 填充关联模型数据
		cols, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to get columns: %v", err)
		}
		vals := make([]interface{}, len(cols))
		for i := range vals {
			vals[i] = new(interface{})
		}

		if err := rows.Scan(vals...); err != nil {
			return fmt.Errorf("failed to scan associated row: %v", err)
		}

		// 将扫描结果映射到结构体字段
		colMap := make(map[string]interface{})
		for i, col := range cols {
			colMap[col] = *(vals[i].(*interface{}))
		}

		// 使用fetcher将数据映射到结构体
		fetcher := &Fetcher{mapValues: convertMapInterfaceToMapString(colMap)}
		if err := e.fetchToStruct(fetcher, elementType, assocModelValue); err != nil {
			return fmt.Errorf("failed to map fetched data to struct: %v", err)
		}

		// 如果原始字段是指针，确保已正确设置
		if fieldType.Kind() == reflect.Ptr && fieldValue.Kind() == reflect.Struct {
			fieldValue.Set(assocModelValue)
		}
	}

	return nil
}

// loadReferenceAssociation 加载引用关联
func (e *Engine) loadReferenceAssociation(fieldValue reflect.Value, fieldType reflect.Type, foreignKey, references string, args ...any) error {
	// 引用关联通常是通过 references 指定的字段进行关联
	// 例如: belongs to 关系，使用 references 指定被引用的字段

	// 获取主模型的引用字段值
	mainModelValue := e.model
	mainModelType := reflect.TypeOf(mainModelValue)
	mainModelReflectValue := reflect.ValueOf(mainModelValue)

	// 处理指针
	if mainModelType.Kind() == reflect.Ptr {
		mainModelType = mainModelType.Elem()
		mainModelReflectValue = mainModelReflectValue.Elem()
	}

	// 检查是否是切片类型，如果是则取第一个元素
	if mainModelType.Kind() == reflect.Slice {
		if mainModelReflectValue.Len() == 0 {
			return fmt.Errorf("main model slice is empty, cannot get reference value")
		}
		// 取第一个元素
		firstElement := mainModelReflectValue.Index(0)
		if firstElement.Kind() == reflect.Ptr {
			if firstElement.IsNil() {
				return fmt.Errorf("first element in main model slice is nil")
			}
			firstElement = firstElement.Elem()
		}
		mainModelType = firstElement.Type()
		mainModelReflectValue = firstElement
	}

	// 查找引用字段的值
	var refValue interface{}
	found := false
	for i := 0; i < mainModelType.NumField(); i++ {
		field := mainModelType.Field(i)
		// 检查字段名是否匹配，或者检查结构体标签中是否有匹配
		if strings.EqualFold(field.Name, references) || strings.EqualFold(e.getTagValue(field), references) {
			fieldVal := mainModelReflectValue.Field(i)
			refValue = fieldVal.Interface()
			found = true
			break
		}
		// 如果没有直接匹配，检查是否有类似名称的字段
		fieldSnakeCase := snakeCase(field.Name)
		if strings.EqualFold(fieldSnakeCase, references) {
			fieldVal := mainModelReflectValue.Field(i)
			refValue = fieldVal.Interface()
			found = true
			break
		}
	}

	if !found {
		// 再次遍历所有字段，尝试找到可能的引用字段
		for i := 0; i < mainModelType.NumField(); i++ {
			field := mainModelType.Field(i)
			// 检查是否有 gorm 标签包含 references
			gormTag := field.Tag.Get("gorm")
			if gormTag != "" {
				settings := parseTagSetting(gormTag, ";")
				// 检查是否在其他字段的 gorm 标签中定义了 references
				if fieldReferences, ok := settings["REFERENCES"]; ok {
					if strings.EqualFold(fieldReferences, references) {
						fieldVal := mainModelReflectValue.Field(i)
						refValue = fieldVal.Interface()
						found = true
						break
					}
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("reference field '%s' value not found in main model", references)
	}

	// 准备关联模型
	var assocModelValue reflect.Value
	if fieldType.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			assocModelValue = reflect.New(fieldType.Elem()).Elem()
			fieldValue.Set(reflect.New(fieldType.Elem()))
			fieldValue = fieldValue.Elem()
		} else {
			assocModelValue = fieldValue.Elem()
		}
	} else {
		assocModelValue = fieldValue
	}

	// 构建查询：SELECT {fields} FROM assoc_table WHERE foreign_key_field = ref_value
	elementType := fieldType
	if fieldType.Kind() == reflect.Ptr {
		elementType = fieldType.Elem()
	}

	tableName := getTableNameFromType(elementType)

	// 动态构建选择字段列表，而不是使用 *
	columns := getStructFields(elementType, e)
	selectFields := strings.Join(columns, ", ")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", selectFields, tableName, foreignKey)

	// 处理不同数据库的占位符
	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", selectFields, tableName, foreignKey)
	}

	rows, err := e.db.Query(query, refValue)
	if err != nil {
		return fmt.Errorf("failed to query reference association: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		// 填充关联模型数据
		cols, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to get columns: %v", err)
		}
		vals := make([]interface{}, len(cols))
		for i := range vals {
			vals[i] = new(interface{})
		}

		if err := rows.Scan(vals...); err != nil {
			return fmt.Errorf("failed to scan associated row: %v", err)
		}

		// 将扫描结果映射到结构体字段
		colMap := make(map[string]interface{})
		for i, col := range cols {
			colMap[col] = *(vals[i].(*interface{}))
		}

		// 使用fetcher将数据映射到结构体
		fetcher := &Fetcher{mapValues: convertMapInterfaceToMapString(colMap)}
		if err := e.fetchToStruct(fetcher, elementType, assocModelValue); err != nil {
			return fmt.Errorf("failed to map fetched data to struct: %v", err)
		}

		// 如果原始字段是指针，确保已正确设置
		if fieldType.Kind() == reflect.Ptr && fieldValue.Kind() == reflect.Struct {
			fieldValue.Set(assocModelValue)
		}
	}

	return nil
}

// getTableNameFromType 从类型获取表名
func getTableNameFromType(t reflect.Type) string {
	// 尝试调用TableName方法
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 创建一个临时实例来调用TableName方法
	tempInstance := reflect.New(t).Interface()
	if tabler, ok := tempInstance.(interface{ TableName() string }); ok {
		return tabler.TableName()
	}

	// 如果没有TableName方法，使用蛇形命名规则
	return snakeCase(t.Name())
}

// snakeCase 转换为蛇形命名
func snakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return strings.ToLower(string(result))
}

// convertMapInterfaceToMapString 将 map[string]interface{} 转换为 map[string]string
func convertMapInterfaceToMapString(input map[string]interface{}) map[string]string {
	result := make(map[string]string, len(input))
	for k, v := range input {
		if v == nil {
			result[k] = ""
		} else {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result
}

// getStructFields 获取结构体的所有数据库字段名
func getStructFields(structType reflect.Type, e *Engine) []string {
	var columns []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// 跳过非导出字段
		if field.PkgPath != "" {
			continue
		}

		// 检查是否是嵌套结构体
		if field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{}) {
			// 递归处理嵌套结构体的字段
			nestedColumns := getStructFields(field.Type, e)
			columns = append(columns, nestedColumns...)
			continue
		}

		// 获取字段的数据库标签
		tagName := e.getTagValue(field)
		if tagName == "" {
			// 如果没有db标签，使用字段名的蛇形格式
			tagName = snakeCase(field.Name)
		}

		// 如果标签值为"-"表示忽略此字段
		if tagName == "-" {
			continue
		}

		columns = append(columns, tagName)
	}

	return columns
}

// isInOperator 检查字符串是否为常见SQL操作符
func isInOperator(op string) bool {
	switch strings.ToUpper(op) {
	case "=", "!=", "<>", "<", "<=", ">", ">=", "LIKE", "NOT LIKE", "IN", "NOT IN", "EXISTS":
		return true
	default:
		return false
	}
}
